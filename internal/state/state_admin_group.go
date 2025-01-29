package state

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/server"
	callback "github.com/epistax1s/gomer/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AdminGroupState struct {
	data *StateContext
}

func NewAdminGroupState(data *StateContext) State {
	return &AdminGroupState{
		data: data,
	}
}

const (
	PAGE_SIZE    = 5
	BTN_PREV     = "<"
	BTN_NEXT     = ">"
	PREV_PREFIX  = "prev_"
	NEXT_PREFIX  = "next_"
	GROUP_PREFIX = "group_"
	MENU         = "menu"
)

type callbackHandler func(*server.Server, *tgbotapi.CallbackQuery)

// TODO rename
var callbacHandlers = map[string]callbackHandler{
	callback.AGPrev:   prevHandler,
	callback.AGNext:   nextHandler,
	callback.AGSelect: groupSelectHandler,
	callback.AGUnlink: groupUnlinkHandler,
}

func prevHandler(server *server.Server, query *tgbotapi.CallbackQuery) {

}

func nextHandler(server *server.Server, query *tgbotapi.CallbackQuery) {

}

func groupSelectHandler(server *server.Server, query *tgbotapi.CallbackQuery) {

}

func groupUnlinkHandler(server *server.Server, query *tgbotapi.CallbackQuery) {

}

func (state *AdminGroupState) Init(server *server.Server, update *tgbotapi.Update) {
	keyboard, err := generateGroupsKeyboard(server, 1)
	if err != nil {
		log.Error(err.Error())
	}
	server.Gomer.SendMessageWithKeyboard(update.FromChat().ID, "Группы", keyboard)
}

func (state *AdminGroupState) Handle(server *server.Server, update *tgbotapi.Update) {
	if update.CallbackQuery != nil {
		jsonData := update.CallbackQuery.Data
		data, err := callback.Decode(jsonData)
		if err != nil {
			log.Error(err.Error())
		}

		callbacHandlers[]

	} else {
		log.Error("ERROR")
	}
}

func handleNavigationBtn(server *server.Server, callback *tgbotapi.CallbackQuery, prefix string) {
	page, err := parseCallbackInt(callback.Data, prefix)
	if err != nil {
		log.Error(err.Error())
	}

	keyboard, err := generateGroupsKeyboard(server, page)
	if err != nil {
		log.Error(err.Error())
	}

	server.Gomer.EditMessageWithKeyboard(callback.From.ID, callback.Message.MessageID, "Группы", keyboard)
}

func handleSelectGroup(server *server.Server, callback *tgbotapi.CallbackQuery) {
	groupID, err := parseCallbackInt(callback.Data, GROUP_PREFIX)
	log.Info("handleSelectGroup", "groupID", groupID)
	if err != nil {
		log.Error(err.Error())
	}

	var keyboard tgbotapi.InlineKeyboardMarkup

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отвязать группу", "unlink"),
	))

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("К списку групп", "back"),
	))

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("В главное меню", MENU),
	))

	// err
	group, _ := server.GroupService.FindByID(int64(groupID))
	log.Info("group", "group", group)

	server.Gomer.EditMessageWithKeyboard(
		callback.From.ID,
		callback.Message.MessageID,
		"Вы выбрали группу"+group.Title+"\nВы можете удалить группу или покинуть данное меню",
		&keyboard)
}

func handleMenuBtn(server *server.Server, update *tgbotapi.Update) {
	StateMachine.
		Set(Idle, update.FromChat().ID, &StateContext{}).
		Init(server, update)
}

func generateGroupsKeyboard(server *server.Server, page int) (*tgbotapi.InlineKeyboardMarkup, error) {
	groups, err := server.GroupService.FindPaginated(page, PAGE_SIZE)
	if err != nil {
		return &tgbotapi.InlineKeyboardMarkup{}, err
	}

	groupCount, err := server.GroupService.CountAll()
	if err != nil {
		return &tgbotapi.InlineKeyboardMarkup{}, err
	}

	var keyboard tgbotapi.InlineKeyboardMarkup

	for _, group := range groups {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(group.Title, GROUP_PREFIX+strconv.Itoa(int(group.ID))),
		))
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("В главное меню", MENU)))

	navigation := tgbotapi.NewInlineKeyboardRow()

	if page > 1 {
		navigation = append(navigation, tgbotapi.NewInlineKeyboardButtonData(BTN_PREV, PREV_PREFIX+strconv.Itoa(page-1)))
	}

	if groupCount > int64(page)*PAGE_SIZE {
		navigation = append(navigation, tgbotapi.NewInlineKeyboardButtonData(BTN_NEXT, NEXT_PREFIX+strconv.Itoa(page+1)))
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, navigation)

	return &keyboard, nil
}

func parseCallbackInt(data string, prefix string) (int, error) {
	parts := strings.Split(data, prefix)
	if len(parts) > 1 {
		return strconv.Atoi(parts[1])
	} else {
		return 0, fmt.Errorf("")
	}
}

/* func parseCallbackString(data string, prefix string) (string, error) {
	parts := strings.Split(data, prefix)
	if len(parts) > 1 {
		return parts[1], nil
	} else {
		return "", fmt.Errorf("")
	}
} */
