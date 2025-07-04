package manage_departs

import (
	"github.com/epistax1s/gomer/internal/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const PAGE_SIZE = 10

func (state *ManageDepartsState) Init(update *tgbotapi.Update) {

}

func (state *ManageDepartsState) Handle(update *tgbotapi.Update) {

}

func (state *ManageDepartsState) prevHandler(update *tgbotapi.Update, с callback.Callback)   {}
func (state *ManageDepartsState) nextHandler(update *tgbotapi.Update, с callback.Callback)   {}
func (state *ManageDepartsState) swapHandler(update *tgbotapi.Update, с callback.Callback)   {}
func (state *ManageDepartsState) addHandeler(update *tgbotapi.Update, с callback.Callback)   {}
func (state *ManageDepartsState) selectHandler(update *tgbotapi.Update, с callback.Callback) {}
func (state *ManageDepartsState) deleteHandler(update *tgbotapi.Update, с callback.Callback) {}
func (state *ManageDepartsState) listHandler(update *tgbotapi.Update, с callback.Callback)   {}
func (state *ManageDepartsState) exitHandler(update *tgbotapi.Update, с callback.Callback)   {}
