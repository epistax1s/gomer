package interceptor

import (
	"github.com/epistax1s/gomer/internal/server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Interceptor interface {
	SetNext(Interceptor)
	Next(*server.Server, *tgbotapi.Update)
	Handle(*server.Server, *tgbotapi.Update)
}

type BaseInterceptor struct {
	next Interceptor
}

func (interceptor *BaseInterceptor) SetNext(next Interceptor) {
	interceptor.next = next
}

func (interceptor *BaseInterceptor) Next(server *server.Server, update *tgbotapi.Update) {
    if interceptor.next != nil {
		interceptor.next.Handle(server, update)
    }
}