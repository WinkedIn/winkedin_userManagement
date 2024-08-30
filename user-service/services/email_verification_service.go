package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/go-redis/redis/v8"
	"github.com/winkedin/user-service/interfaces"
)

const maxRetries = 3

type EmailVerificationServiceImpl struct {
	rdb      *redis.Client
	verifier *emailverifier.Verifier
}

func NewEmailVerificationService(rdb *redis.Client) interfaces.EmailVerificationService {
	v := emailverifier.NewVerifier()
	v = v.EnableSMTPCheck()
	return &EmailVerificationServiceImpl{rdb: rdb, verifier: v}
}

func (e *EmailVerificationServiceImpl) generateOTP(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)[:length]
}

func (e *EmailVerificationServiceImpl) SendVerificationEmail(ctx context.Context, email string) (string, error) {
	otp := e.generateOTP(6)
	otpKey := fmt.Sprintf("otp:%s", email)
	retriesKey := fmt.Sprintf("retries:%s", email)
	err := e.rdb.Set(ctx, otpKey, otp, 10*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store OTP in Redis: %v", err)
	}
	e.rdb.Set(ctx, retriesKey, 0, 10*time.Minute)
	fmt.Printf("Generated OTP for %s: %s\n", email, otp)
	//logic for send email from winkedin email
	return otp, nil
}

func (e *EmailVerificationServiceImpl) VerifyOTP(ctx context.Context, email string, otp string) error {
	otpKey := fmt.Sprintf("otp:%s", email)
	retriesKey := fmt.Sprintf("retries:%s", email)
	storedOTP, err := e.rdb.Get(ctx, otpKey).Result()
	if err == redis.Nil || storedOTP != otp {
		retryCount, _ := e.rdb.Incr(ctx, retriesKey).Result()
		if retryCount >= maxRetries {
			e.rdb.Del(ctx, otpKey)
			return fmt.Errorf("maximum retries reached, OTP invalidated")
		}
		return fmt.Errorf("invalid OTP, retry attempt %d of %d", retryCount, maxRetries)
	}
	e.rdb.Del(ctx, retriesKey)
	return nil
}

func (e *EmailVerificationServiceImpl) ValidateWorkEmail(ctx context.Context, email string) (bool, error) {
	result, err := e.verifier.Verify(email)
	if err != nil {
		return false, fmt.Errorf("failed to verify email: %v", err)
	}
	if result.Disposable {
		// This is a disposable email
		return false, fmt.Errorf("disposable email addresses are not allowed")
	}
	if result.Free {
		// This is a free email service like Gmail, Yahoo, etc.
		return false, fmt.Errorf("invalid email and cannot be used")
	}
	if result.SMTP.Deliverable || result.SMTP.HostExists || result.SMTP.FullInbox || result.SMTP.Disabled {
		// Email address failed SMTP check, which indicates it may not exist
		return false, fmt.Errorf("email address failed verification")
	}
	return true, nil
}
