package user

import (
	"context"
	"errors"
	"github.com/winkedin/user-service/logger"
	"github.com/winkedin/user-service/models"
	"github.com/winkedin/user-service/types"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &handler{
		db: db,
	}
}

// GetUserByEmail - Get user by email
func (h *handler) GetUserByEmail(ctx context.Context, email string) (u *models.User, userExists bool, err error) {
	// log function entrypoint
	logger.InfoLogger.Printf("function-entry: GetUserByEmail, ctx: %v, email: %v", ctx, email)

	var user models.User
	err = h.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.InfoLogger.Printf("function-exit: GetUserByEmail, user: %v, userExists: %v, err: %v", nil, false, nil)
			return nil, false, nil
		}
		logger.ErrorLogger.Printf("function-exit: GetUserByEmail, user: %v, userExists: %v, err: %v", nil, false, err)
		return nil, false, err
	}

	logger.InfoLogger.Printf("function-exit: GetUserByEmail, user: %v, userExists: %v, err: %v", user, true, err)
	return &user, true, err
}

func (h *handler) CreateUserFromLinkedInProfile(ctx context.Context, linkedInUser *types.LinkedInProfile, emailAddress string) (*models.User, error) {
	// log function entrypoint
	logger.InfoLogger.Printf("function-entry: CreateUserFromLinkedInProfile, ctx: %v, linkedInUser: %v", ctx, linkedInUser)

	user := &models.User{
		FirstName:      linkedInUser.LocalizedFirstName,
		LastName:       linkedInUser.LocalizedLastName,
		Email:          emailAddress,
		ProfilePicture: linkedInUser.ProfilePicture.DisplayImage,
	}

	err := h.db.Create(user).Error
	if err != nil {
		logger.ErrorLogger.Printf("function-exit: CreateUserFromLinkedInProfile, user: %v, err: %v", nil, err)
		return nil, err
	}

	logger.InfoLogger.Printf("function-exit: CreateUserFromLinkedInProfile, user: %v, err: %v", user, err)
	return user, nil
}

func (h *handler) UpdateUser(ctx context.Context, user *models.User) error {
	// log function entrypoint
	logger.InfoLogger.Printf("function-entry: UpdateUser, ctx: %v, user: %v", ctx, user)

	err := h.db.Save(user).Error
	if err != nil {
		logger.ErrorLogger.Printf("function-exit: UpdateUser, err: %v", err)
		return err
	}

	logger.InfoLogger.Printf("function-exit: UpdateUser, err: %v", err)
	return nil
}
