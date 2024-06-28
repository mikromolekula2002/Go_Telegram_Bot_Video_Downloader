package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/mikromolekula2002/tg_downloader_bot/pkg/telegram"
)

func main() {
	// достаем токен из .env файла
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки .env файла: ", err)
	}

	token := os.Getenv("Token")
	ApiKey := os.Getenv("ApiKey")
	// инициализация бота и введение его токена
	bot, err := tgbotapi.NewBotAPI(token) //добавить может быть токен через флаг, пока хуй его знает
	if err != nil {
		log.Panic("Ошибка запуска тг бота", err)
	}

	bot.Debug = false // это неудобно, лучше оставлять выкл.

	telegramBot := telegram.NewBot(bot, ApiKey)
	telegramBot.Start()
}
