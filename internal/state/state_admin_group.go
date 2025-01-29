package state

import (
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
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
	PAGE_SIZE = 5
)

type agCallbackHandler func(*server.Server, *tgbotapi.Update, callback.Callback)

var agCallbacHandlers = map[string]agCallbackHandler{
	callback.AG_PREV:       prevHandler,
	callback.AG_NEXT:       nextHandler,
	callback.AG_SELECT:     groupSelectHandler,
	callback.AG_UNLINK:     groupUnlinkHandler,
	callback.AG_GROUP_LIST: groupListHandler,
	callback.EXIT:          exitHandler,
}

func (state *AdminGroupState) Init(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	keyboard, err := genGroupListKeyboard(server, 1)
	if err != nil {
		log.Error("group management / init - group list rendering errors",
			"chatID", chatID, "page", err.Error())

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	server.Gomer.SendMessageWithKeyboard(update.FromChat().ID, i18n.Localize("adminGroupManagementTitle"), keyboard)
}

func (state *AdminGroupState) Handle(server *server.Server, update *tgbotapi.Update) {
	chatID := update.FromChat().ID
	query := update.CallbackQuery

	if query != nil {
		queryData := query.Data
		callback, err := callback.Decode(queryData)
		if err != nil {
			log.Error(err.Error())
		}
		agCallbacHandlers[callback.GetType()](server, update, callback)
	} else {
		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))
	}
}

func prevHandler(server *server.Server, update *tgbotapi.Update, c callback.Callback) {
	prevCallback := c.(*callback.AGPrevCallback)
	page := prevCallback.Page

	renderGroupListKeyboard(server, update, page)
}

func nextHandler(server *server.Server, update *tgbotapi.Update, c callback.Callback) {
	nextCallback := c.(*callback.AGNextCallback)
	page := nextCallback.Page

	renderGroupListKeyboard(server, update, page)
}

func groupSelectHandler(server *server.Server, update *tgbotapi.Update, c callback.Callback) {
	selectCallback := c.(*callback.AGSelectCallback)

	page := selectCallback.Page
	groupID := selectCallback.GroupID

	chatID := update.CallbackQuery.From.ID

	group, err := server.GroupService.FindByID(int64(groupID))
	if err != nil {
		log.Error("group management, error handling the groupSelect button, group not found",
			"groupID", groupID, "callback", c, "err", err.Error())

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	renderGroupControlKeyboard(server, update, page, group)
}

func groupUnlinkHandler(server *server.Server, update *tgbotapi.Update, c callback.Callback) {
	// TODO think about and implement group unlinking
	log.Info("groupUnlinkHandler dummy ", "callback", c)
}

func groupListHandler(server *server.Server, update *tgbotapi.Update, c callback.Callback) {
	groupListCallback := c.(*callback.AGGroupListCallback)
	page := groupListCallback.Page

	renderGroupListKeyboard(server, update, page)
}

func exitHandler(server *server.Server, update *tgbotapi.Update, callback callback.Callback) {
	chatID := update.CallbackQuery.From.ID
	StateMachine.
		Set(Idle, chatID, &StateContext{}).
		Init(server, update)

	server.Gomer.RemoveMarkup(update.CallbackQuery)
}

func renderGroupListKeyboard(server *server.Server, update *tgbotapi.Update, page int) {
	chatID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	keyboard, err := genGroupListKeyboard(server, page)
	if err != nil {
		log.Error("group management - group list rendering errors",
			"chatID", chatID, "page", page, "err", err.Error())

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	server.Gomer.EditMessageWithKeyboard(chatID, messageID, i18n.Localize("adminGroupManagementTitle"), keyboard)
}

func renderGroupControlKeyboard(server *server.Server, update *tgbotapi.Update, page int, group *model.Group) {
	chatID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	keyboard, err := genGroupControlKeyboard(page, int(group.ID))
	if err != nil {
		log.Error("group management - group control rendering errors",
			"chatID", chatID, "page", page, "groupID", int(group.ID), "err", err.Error())

		server.Gomer.SendMessage(chatID, i18n.Localize("oops"))
		return
	}

	server.Gomer.EditMessageWithKeyboard(chatID, messageID, group.Title, keyboard)
}

func genGroupListKeyboard(server *server.Server, page int) (*tgbotapi.InlineKeyboardMarkup, error) {
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
		callbackData := callback.NewAGSelectCallback(page, int(group.ID))
		keyboard.InlineKeyboard = append(
			keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(group.Title, callbackData),
			),
		)
	}

	exitCallbackData := callback.NewExitCallback()
	keyboard.InlineKeyboard = append(
		keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(i18n.Localize("exit"), exitCallbackData),
		),
	)

	navigation := tgbotapi.NewInlineKeyboardRow()
	if page > 1 {
		callbackData := callback.NewAGPrevCallback(page - 1)
		navigation = append(
			navigation,
			tgbotapi.NewInlineKeyboardButtonData(i18n.Localize("adminGroupManagementPrev"), callbackData))
	}

	if groupCount > int64(page)*PAGE_SIZE {
		callbackData := callback.NewAGNextCallback(page + 1)
		navigation = append(
			navigation,
			tgbotapi.NewInlineKeyboardButtonData(i18n.Localize("adminGroupManagementNext"), callbackData))
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, navigation)

	return &keyboard, nil
}

func genGroupControlKeyboard(page int, groupID int) (*tgbotapi.InlineKeyboardMarkup, error) {
	var keyboard tgbotapi.InlineKeyboardMarkup

	unlinkCallbackData := callback.NewAGUnlinkCallback(groupID)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(i18n.Localize("adminGroupManagementUnlink"), unlinkCallbackData),
	))

	backCallbackData := callback.NewAGGroupListCallback(page)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(i18n.Localize("adminGroupManagementGroupList"), backCallbackData),
	))

	exitCallbackData := callback.NewExitCallback()
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(i18n.Localize("exit"), exitCallbackData),
	))

	return &keyboard, nil
}
