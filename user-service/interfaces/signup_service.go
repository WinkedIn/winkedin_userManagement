package interfaces

import (
	"context"

	"github.com/winkedin/user-service/models"
)

type SignupService interface {
	Signup(ctx context.Context, user *models.User) error
}
