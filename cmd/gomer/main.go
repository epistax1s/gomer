package main

import (
	"github.com/epistax1s/gomer/internal/handler"
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

	for update := range updateChan {
		handler.HandleUpdate(server, &update)
	}
}
