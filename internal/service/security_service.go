package service

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/epistax1s/gomer/internal/log"
	"github.com/epistax1s/gomer/internal/model"
)

type SecurityService interface {
	IsRegistered(chatID int64) bool
	IsActive(chatID int64) bool
	IsDeleted(chatID int64) bool
	IsAdmin(chatID int64) bool
	RegisterUser(chatID int64, username string, name string, inviteCode string) error
	AuthenticateUser(chatID int64) (bool, error)
	GenerateInviteLink(creatorChatID int64) (string, error)
	GetUserInvites(chatID int64) ([]model.Invitation, error)
}

type securityService struct {
	userService       UserService
	invitationService InvitationService
	db                *gorm.DB
}

func NewSecurityService(userService UserService, invitationService InvitationService, db *gorm.DB) SecurityService {
	return &securityService{
		userService:       userService,
		invitationService: invitationService,
		db:                db,
	}
}

/*
Checks if a user is registered in the bot.
The presence of a user in the database means that the user is registered in the system.
*/
func (service *securityService) IsRegistered(chatID int64) bool {
	user, err := service.userService.FindByChatID(chatID)
	if err != nil {
		log.Error(
			"SecurityService#IsRegistered() Error loading users",
			"chatID", chatID, "err", err)
		return false
	}

	return user != nil
}

func (service *securityService) IsActive(chatID int64) bool {
	user, err := service.userService.FindByChatID(chatID)
	if err != nil {
		log.Error(
			"SecurityService#IsActive() Error loading users",
			"chatID", chatID, "err", err)
		return false
	}

	if user == nil {
		return false
	}

	return user.Status == model.UserStatusActive
}

func (service *securityService) IsDeleted(chatID int64) bool {
	user, err := service.userService.FindByChatID(chatID)
	if err != nil {
		log.Error(
			"SecurityService#IsDeleted() Error loading users",
			"chatID", chatID, "err", err)
		return false
	}

	if user == nil {
		return false
	}

	return user.Status == model.UserStatusDeleted
}

func (service *securityService) IsAdmin(chatID int64) bool {
	user, err := service.userService.FindByChatID(chatID)
	if err != nil {
		log.Error(
			"SecurityService#IsAdmin() Error loading users",
			"chatID", chatID, "err", err)
		return false
	}

	if user == nil {
		return false
	}

	return user.Role == model.UserRoleAdmin
}

/*
Registers a new user in limbo status.
 1. Validates the invite code
 2. Creates a new user with the given chatID, username, name, and status
 3. Returns an error if the user already exists
*/
func (service *securityService) RegisterUser(chatID int64, username string, name string, inviteCode string) error {
	// Step 1: Validate invitation
	if valid := service.invitationService.ValidateInvite(inviteCode); !valid {
		return fmt.Errorf("invalid invitation code")
	}

	// Step 2: Create new user
	user := &model.User{
		ChatID:    chatID,
		Username:  username,
		Name:      name,
		Status:    model.UserStatusLimbo,
		Role:      model.UserRoleUser,
		CommitSrc: model.UserCommitSrcManual,
		EE:        false,
	}

	if err := service.userService.Create(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Step 3: Mark invitation as used
	if err := service.invitationService.UseInvitation(inviteCode, user); err != nil {
		// TODO: use transaction to rollback user creation if marking invitation fails
		return fmt.Errorf("failed to mark invitation as used: %w", err)
	}

	return nil
}

/*
Authenticates a user by checking if the user exists in the database.
*/
func (service *securityService) AuthenticateUser(chatID int64) (bool, error) {
	user, err := service.userService.FindByChatID(chatID)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

/*
Generates an invite link for a user. The link is used to register a new user in the system.
The link can be used only once.
*/
func (service *securityService) GenerateInviteLink(creatorChatID int64) (string, error) {
	// First check if the user is authenticated
	isAuth, err := service.AuthenticateUser(creatorChatID)
	if err != nil {
		return "", err
	}
	if !isAuth {
		return "", fmt.Errorf("user not authenticated")
	}

	// Generate the invite code
	code, err := service.invitationService.GenerateInvite(creatorChatID)
	if err != nil {
		return "", err
	}

	// Return as a Telegram deep link
	return fmt.Sprintf("https://t.me/your_bot_username?start=%s", code), nil
}

func (service *securityService) GetUserInvites(chatID int64) ([]model.Invitation, error) {
	// First check if the user is authenticated
	isAuth, err := service.AuthenticateUser(chatID)
	if err != nil {
		return nil, err
	}
	if !isAuth {
		return nil, fmt.Errorf("user not authenticated")
	}

	return service.invitationService.GetInvitesByCreator(chatID)
}
