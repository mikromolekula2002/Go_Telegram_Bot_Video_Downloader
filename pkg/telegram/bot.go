package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Структура которая представляет нашего ТГ БОТА
type Bot struct {
	bot    *tgbotapi.BotAPI
	ApiKey string
}

// Инициализирует нашу структуру в мейн
func NewBot(bot *tgbotapi.BotAPI, Apikey string) *Bot {
	return &Bot{
		bot:    bot,
		ApiKey: Apikey,
	}
}

// Самая главная суть в которой у нас открывается бесконечный цикл и раз в 60 сек парсим апдейт(написали ли нашему боту что-нибудь)
func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return fmt.Errorf("\nbot.go - Start.\nОшибка запуска бота и инициализации UpdatesChannel: %w", err)
	}

	b.handleUpdates(updates)
	return nil
}

// отвечает за пришедшие сообщения, команды, нажатия на кнопки
func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.CallbackQuery != nil {
			if err := b.handleCallback(update.CallbackQuery); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
		} else if update.Message != nil {
			if update.Message.IsCommand() {
				err := b.handleCommand(update.Message)
				if err != nil {
					b.handleError(update.Message.Chat.ID, err)
				}
				continue
			} else {
				err := b.handleMessage(update.Message)
				if err != nil {
					b.handleError(update.Message.Chat.ID, err)
				}
			}
		}
	}
}

// создает непосредственный канал который настроен на обновление каждые 60 сек
func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	//установка таймаута в 60 сек, т. е. раз в 60 сек проверяем
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u), nil
}
