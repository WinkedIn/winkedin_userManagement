package interfaces

import "context"

type EmailVerificationService interface {
	SendVerificationEmail(ctx context.Context, email string) (string, error)
	VerifyOTP(ctx context.Context, email string, otp string) error
	ValidateWorkEmail(ctx context.Context, email string) (bool, error)
}
