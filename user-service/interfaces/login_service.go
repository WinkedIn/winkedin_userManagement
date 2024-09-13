package interfaces

import "context"

type LoginService interface {
	Login(ctx context.Context, email string, linkedInJWT string) (string, error)
	Logout(ctx context.Context, userID string) error
}
