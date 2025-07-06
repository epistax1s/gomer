package depart

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"

	. "github.com/epistax1s/gomer/internal/statemachine/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (state *TrackDepartmentState) Init(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	departService := state.server.DepartService

	chatID := update.FromChat().ID

	// get a list of departments
	departments, err := departService.FindAll()
	if err != nil {
		log.Error(
			"error when receiving all departments",
			"chatID", chatID, "state", TrackDepartment, "step", "Init", "err", err)

		gomer.SendMessage(chatID, i18n.Localize("oops"))

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(update)
	}

	if len(departments) == 0 {
		gomer.SendMessage(chatID, i18n.Localize("chooseDepartmentNotFound"))

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(update)
	}

	// creating an inline keyboard for selecting a department
	var departmentMarkup tgbotapi.InlineKeyboardMarkup
	for _, dep := range departments {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(dep.Name, composeDepartmentData(&dep)),
		)
		departmentMarkup.InlineKeyboard = append(departmentMarkup.InlineKeyboard, row)
	}

	// send telegram message
	gomer.SendMessageWithKeyboard(chatID, i18n.Localize("chooseDepartmentPromt"), &departmentMarkup)
}

func (state *TrackDepartmentState) Handle(update *tgbotapi.Update) {
	gomer := state.server.Gomer
	departService := state.server.DepartService

	chatID := update.FromChat().ID
	callbackQuery := update.CallbackQuery

	if callbackQuery == nil {
		log.Warn(
			"callbackQuery is nil",
			"chatID", chatID, "state", TrackDepartment, "step", "Handle")

		state.Init(update)
		return
	}

	callbackData := callbackQuery.Data
	departmentId, extractDepartIdErr := extractDepartmentId(callbackData)
	if extractDepartIdErr != nil {
		log.Error(
			"callbackQuery#data contains an invalid value",
			"chatID", chatID, "state", TrackDepartment, "step", "Handle", "callbackQuery#data", callbackData, "err", extractDepartIdErr)

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(update)

		return
	}

	callbackRespErr := gomer.SendCallbackResponse(callbackQuery, i18n.Localize("chooseDepartmentSuccess"))
	if callbackRespErr != nil {
		log.Error(
			"error confirming callback processing",
			"chatID", chatID, "state", TrackDepartment, "err", callbackRespErr)

	}

	department, err := departService.FindById(departmentId)
	if err != nil {
		log.Error(
			"error when getting the department from the database",
			"chatID", chatID, "state", TrackDepartment, "step", "Handle", "departmentId", departmentId)

		gomer.SendMessage(chatID, i18n.Localize("chooseDepartmentError"))

		state.stateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(update)

		return
	}

	state.stateMachine.
		Set(TrackName, chatID, &StateContext{Department: department}).
		Init(update)

}

func composeDepartmentData(dep *model.Department) string {
	return fmt.Sprintf("department_%v", dep.ID)
}

func extractDepartmentId(depData string) (int64, error) {
	idStr := strings.Split(depData, "_")[1]
	id, err := strconv.Atoi(idStr)
	return int64(id), err
}
