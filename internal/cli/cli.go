package cli

import (
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/repository"
	"github.com/epistax1s/gomer/internal/service"
	"github.com/epistax1s/gomer/internal/config"
)

type CLI struct {
	InvitationService service.InvitationService
	UserService       service.UserService
}

func NewCLI() *CLI {
	config, err := config.LoadConfig()
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}

	log.InitLogger(&config.Log)
	
	db, err := database.InitDatabase()
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}

	err = database.RunMigrations("./database/gomer.db")
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}

	invitationService := service.NewInvitationService(
		repository.NewInvitationRepository(db))

	userService := service.NewUserService(
		repository.NewUserRepository(db))

	return &CLI{
		InvitationService: invitationService,
		UserService:       userService,
	}
	
}

