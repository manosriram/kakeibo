package bot

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/manosriram/kakeibo/internal/handlers"
	"github.com/manosriram/kakeibo/sqlc/db"
)

type TelegramBot struct {
	DB *db.Queries
}

func NewTelegramBot(d *db.Queries) TelegramBot {
	return TelegramBot{
		DB: d,
	}
}

func (t *TelegramBot) HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	var message string
	messageSplits := strings.Split(update.Message.Text, " ")
	if len(messageSplits) > 1 {
		description := strings.Join(messageSplits[1:], " ")

		err := handlers.CreateStatement(t.DB, description)
		if err != nil {
			message = "Error creating statement: " + err.Error()
		} else {
			message = "Statement noted"
		}
	} else {
		message = "Command not found, track expense using /track"
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message,
	})
}

func StartTelegramBot(d *db.Queries) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_ID"), opts...)
	if err != nil {
		panic(err)
	}

	t := NewTelegramBot(d)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/track", bot.MatchTypePrefix, t.HandleMessage)

	b.Start(ctx)
}
