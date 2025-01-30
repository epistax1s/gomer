package interceptor

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Interceptor interface {
	SetNext(Interceptor)
	Next(*tgbotapi.Update)
	Handle(*tgbotapi.Update)
}

type BaseInterceptor struct {
	next Interceptor
}

func (interceptor *BaseInterceptor) SetNext(next Interceptor) {
	interceptor.next = next
}

func (interceptor *BaseInterceptor) Next(update *tgbotapi.Update) {
	if interceptor.next != nil {
		interceptor.next.Handle(update)
	}
}
