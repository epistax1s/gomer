package server

import (
	gomer "github.com/epistax1s/gomer/internal/bot"
	"github.com/epistax1s/gomer/internal/config"
	"github.com/epistax1s/gomer/internal/database"
	"github.com/epistax1s/gomer/internal/i18n"
	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/repository"
	"github.com/epistax1s/gomer/internal/service"
)

type Server struct {
	UserService       service.UserService
	SecurityService   service.SecurityService
	CommitService     service.CommitService
	DepartService     service.DepartService
	GroupService      service.GroupService
	FullCommitService service.FullCommitService
	InvitationService service.InvitationService
	Config            *config.Config
	Gomer             *gomer.Gomer
}

func InitServer() *Server {
	config, err := config.LoadConfig()
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}

	log.InitLogger(&config.Log)
	log.Info("Server initialization: loading the configuration")

	gomer, err := gomer.InitTelegramBot(&config.Bot)
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}
	log.Info("Server initialization: telegram bot initialization")

	db, err := database.InitDatabase()
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}
	log.Info("Server initialization: connecting to the database")

	err = database.RunMigrations("./database/gomer.db")
	if err != nil {
		panic("Server initialization error. err = " + err.Error())
	}
	log.Info("Server initialization: running migrations")

	i18n.InitLocalizer()

	userService := service.NewUserService(
		repository.NewUserRepository(db))

	invitationService := service.NewInvitationService(
		repository.NewInvitationRepository(db))

	securityService := service.NewSecurityService(
		userService,
		invitationService,
		db)

	commitService := service.NewCommitService(
		userService, repository.NewCommitRepository(db))

	departService := service.NewDepartService(
		repository.NewDepartRepository(db))

	groupService := service.NewGroupService(
		repository.NewGroupRepository(db))

	reportService := service.NewFullCommitService(
		repository.NewFullCommit(db))

	return &Server{
		UserService:       userService,
		SecurityService:   securityService,
		CommitService:     commitService,
		DepartService:     departService,
		GroupService:      groupService,
		FullCommitService: reportService,
		InvitationService: invitationService,
		Config:            config,
		Gomer:             gomer,
	}
}
