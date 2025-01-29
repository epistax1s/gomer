package main

import (
	"github.com/epistax1s/gomer/internal/interceptor"
	"github.com/epistax1s/gomer/internal/report"
	"github.com/epistax1s/gomer/internal/server"
	"github.com/epistax1s/gomer/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	server := server.InitServer()

	state.InitStateMachine()

	report.StartNotification(server)
	report.StartPublish(server)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChan := server.Gomer.GetUpdatesChan(updateConfig)

	chain := interceptor.NewChainBuilder().
		Add(&interceptor.LogInterceptor{}).
		Add(&interceptor.RecoverInterceptor{}).
		Add(&interceptor.CancelInterceptor{}).
		Add(&interceptor.HandlerInterceptor{}).
		Build()

	for update := range updateChan {
		chain.Handle(server, &update)
	}
}
