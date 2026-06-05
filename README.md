# Macchina Gialla

Bot Telegram per due utenti: manda una foto di una macchina gialla, aumenta il tuo counter e notifica l'altro.

## Requisiti

- Go 1.22+
- Un bot Telegram (creato via [@BotFather](https://t.me/BotFather))
- Entrambi gli utenti devono aver mandato almeno `/start` al bot prima di ricevere notifiche

## Configurazione

Copia `.env.example` in `.env` e compila i valori:

```env
BOT_TOKEN=      # token del bot da @BotFather
USER_A_ID=      # Telegram user ID del primo utente
USER_B_ID=      # Telegram user ID del secondo utente
USER_A_NAME=    # nome visualizzato nelle notifiche
USER_B_NAME=    # nome visualizzato nelle notifiche
MSG_A=          # messaggio inviato ad A quando manda una foto  (es. "Brava! Totale: {count} 🟡")
MSG_B=          # messaggio inviato a B quando manda una foto
NOTIFY_A=       # messaggio inviato a B quando A manda una foto (es. "A ha avvistato una macchina gialla! Totale: {count} 🟡")
NOTIFY_B=       # messaggio inviato a A quando B manda una foto
COUNTER_FILE=   # percorso del file JSON dei counter (default: counter.json)
```

> Per trovare il proprio user ID: apri [@userinfobot](https://t.me/userinfobot) su Telegram e manda `/start`.

Il placeholder `{count}` nei messaggi viene sostituito con il totale del mittente.

## Avvio

```bash
go run .
```

oppure compila ed esegui:

```bash
go build -o macchina-gialla
./macchina-gialla
```

## Comandi disponibili

| Comando  | Descrizione                              |
|----------|------------------------------------------|
| `/count` | Mostra il totale di entrambi gli utenti  |

## Come funziona

1. Manda una foto al bot
2. Il tuo counter aumenta di 1
3. L'altro utente riceve la foto con il messaggio di notifica configurato
4. Se la foto ha una didascalia, viene inclusa nella notifica

I counter sono salvati in `counter.json` e persistono tra i riavvii.
