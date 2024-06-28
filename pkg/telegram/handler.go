package telegram

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	downloader "github.com/mikromolekula2002/tg_downloader_bot/pkg/Downloader"
)

// метод который проверяет пришла ли нам команда и существует ли она вообще
func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	switch message.Command() {
	case "start":
		// команда START бота
		return b.handleStartCommand(message)
	case "help":
		return b.handleHelpCommand(message)
	default:
		// неизвестная команда
		return b.handleUnknownCommand(message)
	}
}

// метод отвечает за ответ на обычные сообщения отправленные боту(через него будем обрабатывать ссылки на видосы)
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	var videoName string
	var err error

	log.Println("Обработка сообщения пользователя.")
	// ПРОПИСАТЬ СВИТЧ БЛОК НА ПРОВЕРКУ ИСТОЧНИКА ВИДЕО УРЛ
	switch {
	case strings.Contains(message.Text, "youtube.com/shorts"):
		videoName, err = downloader.GetShorts(message.Text, b.ApiKey)
		if err != nil {
			fmt.Println(err)
			return errTikTok
		}
	case strings.Contains(message.Text, "tiktok.com"):
		videoName, err = downloader.GetTiktok(message.Text, b.ApiKey)
		if err != nil {
			fmt.Println(err)
			return errTikTok
		}

	case strings.Contains(message.Text, "instagram.com/reel"):
		videoName, err = downloader.GetReels(message.Text, b.ApiKey)
		if err != nil {
			fmt.Println(err)
			return errReels
		}

	case strings.Contains(message.Text, "youtu.be"):
		return errYoutube
	default:
		return errInvalidURL
	}
	// Открываем видеофайл
	videoFile := tgbotapi.FilePath(videoName)
	// Создаем сообщение с видео
	msg := tgbotapi.NewVideo(message.Chat.ID, videoFile)

	// Отправляем сообщение
	if _, err := b.bot.Send(msg); err != nil {
		fmt.Println("handler.go - handleMessage: ", errSendFailed)
		return err
	}
	mseg := tgbotapi.NewMessage(message.Chat.ID, "Если вы хотите скачать еще видео, просто отправьте нам след-ую ссылку.")
	// Отправляем сообщение
	if _, err := b.bot.Send(mseg); err != nil {
		fmt.Println("handler.go - handleMessage: ", errSendFailed)
		return err
	}
	log.Println("Видео успешно отправлено пользователю.")
	downloader.DeleteVideo(videoName)
	return nil
}

// реализация метода команды START
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	log.Println("Обработка команды '/Start'.\n.")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Бот предназначен для скачивания видео с различных площадок: tiktok, youtube, instagram.")
	// Передаем также INLINE Keyboard с кнопками ознакомления
	msg.ReplyMarkup = createMainMenuKeyboard()
	// Отправляем сообщение
	_, err := b.bot.Send(msg)
	if err != nil {
		fmt.Println("handler.go - handleStartCommand: ", errSendFailed)
		return err
	}
	return nil
}

// реализация метода команды START
func (b *Bot) handleHelpCommand(message *tgbotapi.Message) error {
	log.Println("Обработка команды '/Help'.\n.")
	msg := tgbotapi.NewMessage(message.Chat.ID, "Бот предназначен для скачивания видео с различных площадок.\nТаких как: TikTok, Youtube, Instagram.\nПока что бот работает на бесплатной основе, в скором будущем будет реализована работа подписок.\n Важно: работа с Youtube, а именно НЕ с Youtube Shorts временно невозможна. Приносим свои извинения.\n Для скачивания видео достаточно просто отправить ссылку на него.")
	// Передаем также INLINE Keyboard с кнопками ознакомления
	msg.ReplyMarkup = createMainMenuKeyboard()
	// Отправляем сообщение
	_, err := b.bot.Send(msg)
	if err != nil {
		fmt.Println("handler.go - handleHelpCommand: ", errSendFailed)
		return err
	}
	return nil
}

// реализация метода неизвестной команды
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	log.Println("Обработка неизвестной команды.\n.")
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%s . Такой команды не существует.Ознакомьтесь с работой бота.", message.Text))
	// Передаем также INLINE Keyboard с кнопками ознакомления
	msg.ReplyMarkup = createMainMenuKeyboard()
	_, err := b.bot.Send(msg)
	if err != nil {
		fmt.Println("handleUnkownCommand: Обработан запрос несуществующей команды:", message.Text)
		return err
	}
	return nil
}

// метод который проверяет пришла ли нам команда и существует ли она вообще
func (b *Bot) handleCallback(CallbackQuery *tgbotapi.CallbackQuery) error {
	log.Println("Обработка нажатия на нашу клавиатуру.\n.")

	//Заранее объявляем мисейдж
	var msgString string
	var photoFilePath string

	// Проверка какая кнопка нажата:
	switch CallbackQuery.Data {
	case "shorts":
		msgString = "Бот умеет загружать видео с площадки Youtube Shorts.\n Для работы бота отправьте ссылку на видео, которое хотите скачать."
		photoFilePath = "./Photos/YoutubeShorts.png"

	case "reels":
		msgString = "Бот умеет загружать видео с площадки Instagram Reels.\n Для работы бота отправьте ссылку на видео, которое хотите скачать."
		photoFilePath = "./Photos/Instagram.jpg"

	case "tiktok":
		msgString = "Бот умеет загружать видео с площадки Tiktok.\n Для работы бота отправьте ссылку на видео, которое хотите скачать."
		photoFilePath = "./Photos/TikTok.png"

	case "subscribe":
		msgString = "Для работы бота требуется приобретение подписки, ознакомтесь со следующим:"
		photoFilePath = "./Photos/Subscribe.jpg"

	}

	photoFile, err := os.Open(photoFilePath)
	if err != nil {
		log.Fatalf("failed to open photo file: %v", err)
	}
	defer photoFile.Close()

	photo := tgbotapi.NewPhoto(CallbackQuery.Message.Chat.ID, tgbotapi.FileReader{
		Name:   photoFilePath,
		Reader: photoFile,
	})

	photo.Caption = msgString
	photo.ReplyMarkup = createMainMenuKeyboard()

	// Отправляем сообщение с фото
	if _, err := b.bot.Send(photo); err != nil {
		fmt.Println("handler.go - handleCallback: ", errSendFailed)
		return err
	}

	return nil
}
