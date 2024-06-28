package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func createMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Youtube Shorts", "shorts"),
			tgbotapi.NewInlineKeyboardButtonData("Instagram Reels", "reels"),
			tgbotapi.NewInlineKeyboardButtonData("Tiktok", "tiktok"),
			tgbotapi.NewInlineKeyboardButtonData("Подписка", "subscribe"),
		),
	)
	return inlineKeyboard
}
