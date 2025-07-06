package interceptor

import (
	"fmt"

	"github.com/epistax1s/gomer/internal/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LogInterceptor struct {
	BaseInterceptor
}

func (i *LogInterceptor) Handle(update *tgbotapi.Update) {
	// Determine update type
	updateType := i.getUpdateType(update)

	// Log basic update information
	log.Debug("[LogInterceptor] Update received",
		"updateID", update.UpdateID,
		"type", updateType,
	)

	// Log detailed information based on update type
	switch updateType {
	case "message":
		i.logMessage(update.Message)
	case "callback_query":
		i.logCallbackQuery(update.CallbackQuery)
	case "edited_message":
		i.logMessage(update.EditedMessage)
	case "channel_post":
		i.logMessage(update.ChannelPost)
	case "edited_channel_post":
		i.logMessage(update.EditedChannelPost)
	case "inline_query":
		i.logInlineQuery(update.InlineQuery)
	case "chosen_inline_result":
		i.logChosenInlineResult(update.ChosenInlineResult)
	case "shipping_query":
		i.logShippingQuery(update.ShippingQuery)
	case "pre_checkout_query":
		i.logPreCheckoutQuery(update.PreCheckoutQuery)
	case "poll":
		i.logPoll(update.Poll)
	case "poll_answer":
		i.logPollAnswer(update.PollAnswer)
	case "my_chat_member":
		i.logChatMemberUpdated(update.MyChatMember)
	case "chat_member":
		i.logChatMemberUpdated(update.ChatMember)
	case "chat_join_request":
		i.logChatJoinRequest(update.ChatJoinRequest)
	default:
		log.Debug("[LogInterceptor] Unknown update type")
	}

	i.Next(update)
}

func (i *LogInterceptor) getUpdateType(update *tgbotapi.Update) string {
	if update.Message != nil {
		return "message"
	}
	if update.CallbackQuery != nil {
		return "callback_query"
	}
	if update.EditedMessage != nil {
		return "edited_message"
	}
	if update.ChannelPost != nil {
		return "channel_post"
	}
	if update.EditedChannelPost != nil {
		return "edited_channel_post"
	}
	if update.InlineQuery != nil {
		return "inline_query"
	}
	if update.ChosenInlineResult != nil {
		return "chosen_inline_result"
	}
	if update.ShippingQuery != nil {
		return "shipping_query"
	}
	if update.PreCheckoutQuery != nil {
		return "pre_checkout_query"
	}
	if update.Poll != nil {
		return "poll"
	}
	if update.PollAnswer != nil {
		return "poll_answer"
	}
	if update.MyChatMember != nil {
		return "my_chat_member"
	}
	if update.ChatMember != nil {
		return "chat_member"
	}
	if update.ChatJoinRequest != nil {
		return "chat_join_request"
	}
	return "unknown"
}

// formatUserInfo formats user information for logging
func (i *LogInterceptor) formatUserInfo(user interface{}) string {
	switch u := user.(type) {
	case *tgbotapi.User:
		if u == nil {
			return "[unknown]"
		}
		return fmt.Sprintf("[username: %s, userID: %d]",
			u.UserName,
			u.ID)
	case tgbotapi.User:
		return fmt.Sprintf("[username: %s, userID: %d]",
			u.UserName,
			u.ID)
	default:
		return "[unknown]"
	}
}

// formatChatInfo formats chat information for logging
func (i *LogInterceptor) formatChatInfo(chat interface{}) string {
	switch c := chat.(type) {
	case *tgbotapi.Chat:
		if c == nil {
			return "[unknown]"
		}
		chatTitle := c.Title
		if chatTitle == "" {
			chatTitle = "private"
		}
		return fmt.Sprintf("[title: %s, chatID: %d, type: %s]",
			chatTitle,
			c.ID,
			c.Type)
	case tgbotapi.Chat:
		chatTitle := c.Title
		if chatTitle == "" {
			chatTitle = "private"
		}
		return fmt.Sprintf("[title: %s, chatID: %d, type: %s]",
			chatTitle,
			c.ID,
			c.Type)
	default:
		return "[unknown]"
	}
}

func (i *LogInterceptor) logMessage(message *tgbotapi.Message) {
	if message == nil {
		return
	}

	userInfo := i.formatUserInfo(message.From)
	chatInfo := i.formatChatInfo(message.Chat)

	text := "no text"
	if message.Text != "" {
		text = message.Text
	}

	log.Debug("[LogInterceptor] Message details",
		"user", userInfo,
		"chat", chatInfo,
		"text", text,
		"date", message.Date,
		"messageID", message.MessageID,
	)
}

func (i *LogInterceptor) logCallbackQuery(callback *tgbotapi.CallbackQuery) {
	if callback == nil {
		return
	}

	userInfo := i.formatUserInfo(callback.From)

	messageID := 0
	if callback.Message != nil {
		messageID = callback.Message.MessageID
	}

	log.Debug("[LogInterceptor] CallbackQuery details",
		"user", userInfo,
		"data", callback.Data,
		"chatInstance", callback.ChatInstance,
		"messageID", messageID,
	)
}

func (i *LogInterceptor) logInlineQuery(query *tgbotapi.InlineQuery) {
	if query == nil {
		return
	}

	userInfo := i.formatUserInfo(query.From)

	queryText := query.Query
	if queryText == "" {
		queryText = "[empty query]"
	}

	log.Debug("[LogInterceptor] InlineQuery details",
		"user", userInfo,
		"query", queryText,
		"offset", query.Offset,
	)
}

func (i *LogInterceptor) logChosenInlineResult(result *tgbotapi.ChosenInlineResult) {
	if result == nil {
		return
	}

	userInfo := i.formatUserInfo(result.From)

	resultID := result.ResultID
	if resultID == "" {
		resultID = "[empty resultID]"
	}

	queryText := result.Query
	if queryText == "" {
		queryText = "[empty query]"
	}

	log.Debug("[LogInterceptor] ChosenInlineResult details",
		"user", userInfo,
		"resultID", resultID,
		"query", queryText,
	)
}

func (i *LogInterceptor) logShippingQuery(query *tgbotapi.ShippingQuery) {
	if query == nil {
		return
	}

	userInfo := i.formatUserInfo(query.From)

	invoicePayload := query.InvoicePayload
	if invoicePayload == "" {
		invoicePayload = "[empty payload]"
	}

	log.Debug("[LogInterceptor] ShippingQuery details",
		"user", userInfo,
		"invoicePayload", invoicePayload,
		"shippingAddress", query.ShippingAddress,
	)
}

func (i *LogInterceptor) logPreCheckoutQuery(query *tgbotapi.PreCheckoutQuery) {
	if query == nil {
		return
	}

	userInfo := i.formatUserInfo(query.From)

	invoicePayload := query.InvoicePayload
	if invoicePayload == "" {
		invoicePayload = "[empty payload]"
	}

	currency := query.Currency
	if currency == "" {
		currency = "[empty currency]"
	}

	log.Debug("[LogInterceptor] PreCheckoutQuery details",
		"user", userInfo,
		"invoicePayload", invoicePayload,
		"currency", currency,
		"totalAmount", query.TotalAmount,
	)
}

func (i *LogInterceptor) logPoll(poll *tgbotapi.Poll) {
	if poll == nil {
		return
	}

	log.Debug("[LogInterceptor] Poll details",
		"id", poll.ID,
		"question", poll.Question,
		"optionsCount", len(poll.Options),
		"totalVoterCount", poll.TotalVoterCount,
	)
}

func (i *LogInterceptor) logPollAnswer(answer *tgbotapi.PollAnswer) {
	if answer == nil {
		return
	}

	userInfo := i.formatUserInfo(answer.User)

	log.Debug("[LogInterceptor] PollAnswer details",
		"user", userInfo,
		"pollID", answer.PollID,
		"optionIDs", answer.OptionIDs,
	)
}

func (i *LogInterceptor) logChatMemberUpdated(update *tgbotapi.ChatMemberUpdated) {
	if update == nil {
		return
	}

	userInfo := i.formatUserInfo(update.From)
	chatInfo := i.formatChatInfo(update.Chat)

	log.Debug("[LogInterceptor] ChatMemberUpdated details",
		"user", userInfo,
		"chat", chatInfo,
		"date", update.Date,
		"oldChatMember", update.OldChatMember,
		"newChatMember", update.NewChatMember,
	)
}

func (i *LogInterceptor) logChatJoinRequest(request *tgbotapi.ChatJoinRequest) {
	if request == nil {
		return
	}

	userInfo := i.formatUserInfo(request.From)
	chatInfo := i.formatChatInfo(request.Chat)

	log.Debug("[LogInterceptor] ChatJoinRequest details",
		"user", userInfo,
		"chat", chatInfo,
		"date", request.Date,
		"bio", request.Bio,
		"inviteLink", request.InviteLink,
	)
}
