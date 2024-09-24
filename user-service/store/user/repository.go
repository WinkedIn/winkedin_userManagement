package user

import (
	"context"
	"github.com/winkedin/user-service/models"
	"github.com/winkedin/user-service/types"
)

type UserStore interface {
	GetUserByEmail(context.Context, string) (*models.User, bool, error)
	CreateUserFromLinkedInProfile(ctx context.Context, linkedInUser *types.LinkedInProfile, emailAddress string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}
