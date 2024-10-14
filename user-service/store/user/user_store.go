package user

import (
	"context"
	"errors"
	"github.com/winkedin/user-service/constants"
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
	logger.LogFunctionPointWithContext(ctx, constants.LogFunctionEntry)

	var user models.User
	err = h.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.LogErrorWithContext(ctx, gorm.ErrRecordNotFound.Error())
		return nil, false, nil
	}
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return nil, false, err
	}

	defer logger.LogFunctionPointWithContext(ctx, constants.LogFunctionExit)
	return &user, true, err
}

func (h *handler) CreateUserFromLinkedInProfile(ctx context.Context, linkedInUser *types.LinkedInProfile, emailAddress string) (*models.User, error) {
	// log function entrypoint
	logger.LogFunctionPointWithContext(ctx, constants.LogFunctionEntry)

	user := &models.User{
		FirstName:      linkedInUser.LocalizedFirstName,
		LastName:       linkedInUser.LocalizedLastName,
		Email:          emailAddress,
		ProfilePicture: linkedInUser.ProfilePicture.DisplayImage,
	}

	err := h.db.Create(user).Error
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return nil, err
	}

	logger.LogFunctionPointWithContext(ctx, constants.LogFunctionExit)
	return user, nil
}

func (h *handler) UpdateUser(ctx context.Context, user *models.User) error {
	// log function entrypoint
	logger.LogFunctionPointWithContext(ctx, constants.LogFunctionEntry)

	err := h.db.Save(user).Error
	if err != nil {
		logger.LogErrorWithContext(ctx, err.Error())
		return err
	}

	logger.LogFunctionPointWithContext(ctx, constants.LogFunctionExit)
	return nil
}
