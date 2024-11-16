package state

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/epistax1s/gomer/internal/server"
)

type TrackDepartmentState struct {
	data *StateContext
}

func NewTrackDepartmentState(data *StateContext) State {
	return &TrackDepartmentState{
		data: data,
	}
}

func (state *TrackDepartmentState) Init(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	// get a list of departments
	departments, err := server.DepartService.FindAll()
	if err != nil {
		log.Error(
			"error when receiving all departments",
			"chatID", chatID, "state", TrackDepartment, "step", "Init", "err", err)

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))

		StateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(server, update)
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
	server.Gomer.SendMessageWithMarkup(chatID, i18n.Localize("chooseDepartmentPromt"), departmentMarkup)
}

func (state *TrackDepartmentState) Handle(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID
	callbackQuery := update.CallbackQuery

	if callbackQuery == nil {
		log.Warn(
			"callbackQuery is nil",
			"chatID", chatID, "state", TrackDepartment, "step", "Handle")

		state.Init(server, update)
		return
	}

	callbackData := callbackQuery.Data
	departmentId, extractDepartIdErr := extractDepartmentId(callbackData)
	if extractDepartIdErr != nil {
		log.Error(
			"callbackQuery#data contains an invalid value",
			"chatID", chatID, "state", TrackDepartment, "step", "Handle", "callbackQuery#data", callbackData, "err", extractDepartIdErr)

		StateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(server, update)

		return
	}

	callbackRespErr := server.Gomer.SendCallbackResponse(callbackQuery, i18n.Localize("chooseDepartmentSuccess"))
	if callbackRespErr != nil {
		log.Error(
			"error confirming callback processing",
			"chatID", chatID, "state", TrackDepartment, "err", callbackRespErr)

	}

	department, err := server.DepartService.FindById(departmentId)
	if err != nil {
		log.Error(
			"error when getting the department from the database",
			"chatID", chatID, "state", TrackDepartment, "step", "Handle", "departmentId", departmentId)

		server.Gomer.SendMessage(chatID, i18n.Localize("chooseDepartmentError"))

		StateMachine.
			Set(Idle, chatID, &StateContext{}).
			Init(server, update)

		return
	}

	StateMachine.
		Set(TrackName, chatID, &StateContext{Department: department}).
		Init(server, update)

}

func composeDepartmentData(dep *model.Department) string {
	return fmt.Sprintf("department_%v", dep.ID)
}

func extractDepartmentId(depData string) (int64, error) {
	idStr := strings.Split(depData, "_")[1]
	id, err := strconv.Atoi(idStr)
	return int64(id), err
}
