package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var allowedIDs [2]int64

type counterStore struct {
	CountA int `json:"count_a"`
	CountB int `json:"count_b"`
}

type user struct {
	id        int64
	name      string
	replyMsg  string
	notifyMsg string
}

func isAllowed(userID int64) bool {
	return userID == allowedIDs[0] || userID == allowedIDs[1]
}

func loadCounter(path string) (counterStore, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return counterStore{}, nil
	}
	if err != nil {
		return counterStore{}, err
	}
	var s counterStore
	if err := json.Unmarshal(data, &s); err != nil {
		return counterStore{}, err
	}
	return s, nil
}

func saveCounter(path string, s counterStore) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func formatMsg(template string, count int) string {
	return strings.ReplaceAll(template, "{count}", strconv.Itoa(count))
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return v
}

func mustParseID(key string) int64 {
	v := mustEnv(key)
	id, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return id
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	_ = godotenv.Load()

	userA := user{
		id:        mustParseID("USER_A_ID"),
		name:      mustEnv("USER_A_NAME"),
		replyMsg:  mustEnv("MSG_A"),
		notifyMsg: mustEnv("NOTIFY_A"),
	}
	userB := user{
		id:        mustParseID("USER_B_ID"),
		name:      mustEnv("USER_B_NAME"),
		replyMsg:  mustEnv("MSG_B"),
		notifyMsg: mustEnv("NOTIFY_B"),
	}
	allowedIDs = [2]int64{userA.id, userB.id}

	counterPath := envOr("COUNTER_FILE", "counter.json")
	counts, err := loadCounter(counterPath)
	if err != nil {
		log.Fatalf("load counter: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(mustEnv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("bot init: %v", err)
	}
	log.Printf("bot started as @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := update.Message

		if !isAllowed(msg.From.ID) { // V1
			continue
		}

		var sender, other user
		if msg.From.ID == userA.id {
			sender, other = userA, userB
		} else {
			sender, other = userB, userA
		}

		switch {
		case msg.Photo != nil: // V2, V6
			if sender.id == userA.id {
				counts.CountA++
			} else {
				counts.CountB++
			}
			if err := saveCounter(counterPath, counts); err != nil { // V4
				log.Printf("save counter: %v", err)
			}
			senderCount := counts.CountA
			if sender.id == userB.id {
				senderCount = counts.CountB
			}
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, formatMsg(sender.replyMsg, senderCount)))
			bot.Send(tgbotapi.NewMessage(other.id, formatMsg(sender.notifyMsg, senderCount))) // V3

		case msg.IsCommand() && msg.Command() == "count": // I.cmd
			reply := fmt.Sprintf("%s: %d\n%s: %d", userA.name, counts.CountA, userB.name, counts.CountB)
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, reply))
		}
	}
}
