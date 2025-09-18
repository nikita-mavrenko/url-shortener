package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nikitamavrenko/telegram-bot-service/internal/utils"
)

type Bot struct {
	api       *tgbotapi.BotAPI
	shortener Shortener
}

type Shortener interface {
	Short(ctx context.Context, url string) (string, error)
}

func NewBot(token string, shortener Shortener) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}
	return &Bot{
		api:       api,
		shortener: shortener,
	}, nil
}

func (b *Bot) Start() error {
	return b.startPolling()
}

func (b *Bot) startPolling() error {
	u := tgbotapi.NewUpdate(0)
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if err := b.handle(update); err != nil {
			return fmt.Errorf("failed to handle update: %w", err)
		}
	}

	return nil
}

func (b *Bot) sendMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)

	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (b *Bot) handle(update tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	if update.Message.Text == "/start" {
		return b.handleStart(update.Message.Chat.ID)
	}

	if utils.IsValidLink(update.Message.Text) {
		return b.handleLink(update.Message.Chat.ID, update.Message.Text)
	}

	return nil
}

func (b *Bot) handleStart(id int64) error {
	msg := "Отправь мне любую ссылку и я ее сокращу"
	return b.sendMessage(id, msg)
}

func (b *Bot) handleLink(id int64, link string) error {
	msg := "Ваша ссылка:"

	shortenedLink, err := b.shortener.Short(context.Background(), link)
	if err != nil {
		return fmt.Errorf("failed to shorten link: %w", err)
	}

	return b.sendMessage(id, fmt.Sprintf("%s %s", msg, shortenedLink))
}
