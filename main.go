package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var allowedIDs [2]int64

type counterStore struct {
	Count int `json:"count"`
}

func isAllowed(userID int64) bool {
	return userID == allowedIDs[0] || userID == allowedIDs[1]
}

func otherUser(senderID int64) int64 {
	if senderID == allowedIDs[0] {
		return allowedIDs[1]
	}
	return allowedIDs[0]
}

func loadCounter(path string) (int, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	var s counterStore
	if err := json.Unmarshal(data, &s); err != nil {
		return 0, err
	}
	return s.Count, nil
}

func saveCounter(path string, count int) error {
	data, err := json.Marshal(counterStore{Count: count})
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
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

	allowedIDs[0] = mustParseID("USER_A_ID")
	allowedIDs[1] = mustParseID("USER_B_ID")
	counterPath := envOr("COUNTER_FILE", "counter.json")

	count, err := loadCounter(counterPath)
	if err != nil {
		log.Fatalf("load counter: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(mustEnv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("bot init: %v", err)
	}
	log.Printf("bot started as @%s, counter=%d", bot.Self.UserName, count)

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

		switch {
		case msg.Photo != nil: // V2
			count++
			if err := saveCounter(counterPath, count); err != nil { // V4
				log.Printf("save counter: %v", err)
			}
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Macchina gialla! Totale: %d 🟡", count)))
			bot.Send(tgbotapi.NewMessage(otherUser(msg.From.ID), fmt.Sprintf("Macchina gialla avvistata! Totale: %d 🟡", count))) // V3

		case msg.IsCommand() && msg.Command() == "count": // I.cmd
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Macchine gialle avvistate: %d 🟡", count)))
		}
	}
}
