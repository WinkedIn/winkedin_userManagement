package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/winkedin/user-service/interfaces"
	"github.com/winkedin/user-service/models"
	"gorm.io/gorm"
)

type SignupServiceImpl struct {
	db                   *gorm.DB
	rdb                  *redis.Client
	emailVerificationSvc interfaces.EmailVerificationService
}

func NewSignupService(db *gorm.DB, rdb *redis.Client, emailVerificationSvc interfaces.EmailVerificationService) interfaces.SignupService {
	return &SignupServiceImpl{
		db:                   db,
		rdb:                  rdb,
		emailVerificationSvc: emailVerificationSvc,
	}
}

func (s *SignupServiceImpl) Signup(ctx context.Context, user *models.User) error {
	isValidWorkEmail, err := s.emailVerificationSvc.ValidateWorkEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to validate work email: %v", err)
	}
	if !isValidWorkEmail {
		return fmt.Errorf("provided email is a public or proxy email; please use a work email")
	}
	userKey := fmt.Sprintf("user:%s", user.Email)
	err = s.rdb.HMSet(ctx, userKey, map[string]interface{}{
		"FirstName":      user.FirstName,
		"LastName":       user.LastName,
		"Email":          user.Email,
		"LinkedInID":     user.LinkedInID,
		"CompanyName":    user.CompanyName,
		"JobTitle":       user.JobTitle,
		"ProfilePicture": user.ProfilePicture,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to store user data in Redis: %v", err)
	}
	s.rdb.Expire(ctx, userKey, 15*time.Minute)
	_, err = s.emailVerificationSvc.SendVerificationEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	return nil
}
