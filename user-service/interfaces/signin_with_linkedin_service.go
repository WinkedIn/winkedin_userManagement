package interfaces

import "context"

type SignInWithLinkedInService interface {
	GetLinkedInProfileAndLogin(ctx context.Context, code string) (string, error)
}
