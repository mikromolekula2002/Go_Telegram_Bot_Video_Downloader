package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errSendFailed = errors.New("ошибка отправки пользователю сообщения")
	errYoutube    = errors.New("не рабочий метод скачивания с Youtube, пока в реализации")
	errInvalidURL = errors.New("невалидная ссылка")
	errStandart   = errors.New("произошла ошибка, повторите позже")
	errTikTok     = errors.New("ошибка скачивания TikTok ролика")
	errShorts     = errors.New("ошибка скачивания Youtube Shorts ролика")
	errReels      = errors.New("ошибка скачивания Instagram Reels ролика")
)

func (b *Bot) handleError(chatID int64, err error) {

	msg := tgbotapi.NewMessage(chatID, "Произошла ошибка, попробуйте позже")

	switch err {
	case errInvalidURL:
		msg.Text = "Невалидная ссылка. Пожалуйста ознакомьтесь с какими платформами может работать телеграмм бот."
		b.bot.Send(msg)
	case errStandart:
		msg.Text = "Произошла ошибка, повторите позже"
		b.bot.Send(msg)
	case errTikTok:
		msg.Text = "Произошла ошибка со скачиванием TikTok ролика, попробуйте позже"
		b.bot.Send(msg)
	case errShorts:
		msg.Text = "Произошла ошибка со скачиванием Youtube Shorts ролика, попробуйте позже"
		b.bot.Send(msg)
	case errReels:
		msg.Text = "Произошла ошибка со скачиванием Reels ролика, попробуйте позже"
		b.bot.Send(msg)
	case errYoutube:
		msg.Text = "Извините, пока что телеграмм бот не умеет работать с обычными видеороликами с платформы Youtube. Вы можете скачать видеоролик Youtube Shorts, для этого просто отправьте ссылку на видео."
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}
}
