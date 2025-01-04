package main

import (
	"github.com/epistax1s/gomer/internal/interceptor"
	"github.com/epistax1s/gomer/internal/report"
	"github.com/epistax1s/gomer/internal/server"
	"github.com/epistax1s/gomer/internal/statemachine/builder"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	server := server.InitServer()

	stateMachine := builder.NewStateMachine(server)

	report.StartNotification(server)
	report.StartPublish(server)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updateChan := server.Gomer.GetUpdatesChan(updateConfig)

	chain := interceptor.NewChainBuilder().
		Add(&interceptor.LogInterceptor{}).
		Add(&interceptor.RecoverInterceptor{}).
		Add(&interceptor.CancelInterceptor{
			Server:       server,
			StateMachine: stateMachine,
		}).
		Add(&interceptor.HandlerInterceptor{
			Server:       server,
			StateMachine: stateMachine,
		}).
		Build()

	for update := range updateChan {
		chain.Handle(&update)
	}
}
