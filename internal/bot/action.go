package gomer

import (
	"github.com/epistax1s/gomer/internal/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (gomer *Gomer) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	gomer.Send(msg)
}

func (gomer *Gomer) SendMessageWithMarkup(chatID int64, text string, markup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = markup
	gomer.Send(msg)
}

func (gomer *Gomer) SendCallbackResponse(callback *tgbotapi.CallbackQuery, text string) error {
	callbackConfig := tgbotapi.CallbackConfig{
		CallbackQueryID: callback.ID,
		Text:            text,
		ShowAlert:       false,
	}

	_, err := gomer.Request(callbackConfig)
	return err
}

func (gomer *Gomer) RemoveMarkup(callback *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, callback.Message.Text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)

	// TODO fix log
	if _, err := gomer.Send(msg); err != nil {
		log.Error(
			"error",
			"err", err)
	}
}
