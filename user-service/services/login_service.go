package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/winkedin/user-service/interfaces"
	"github.com/winkedin/user-service/models"
	"gorm.io/gorm"
)

type LoginServiceImpl struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewLoginService(db *gorm.DB, rdb *redis.Client) interfaces.LoginService {
	return &LoginServiceImpl{db: db, rdb: rdb}
}

func (s *LoginServiceImpl) Login(ctx context.Context, email string, linkedInJWT string) (string, error) {
	token, err := jwt.Parse(linkedInJWT, nil)
	if err != nil {
		return "", fmt.Errorf("invalid LinkedIn token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid LinkedIn token claims")
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", fmt.Errorf("failed to get expiration from LinkedIn token")
	}
	expirationTime := time.Unix(int64(exp), 0)
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", fmt.Errorf("user not found")
	}
	tokenKey := fmt.Sprintf("session:%s", user.ID)
	err = s.rdb.Set(ctx, tokenKey, linkedInJWT, expirationTime.Sub(time.Now())).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store session token in Redis: %v", err)
	}
	return linkedInJWT, nil
}

func (s *LoginServiceImpl) Logout(ctx context.Context, userID string) error {
	tokenKey := fmt.Sprintf("session:%s", userID)
	err := s.rdb.Del(ctx, tokenKey).Err()
	if err != nil {
		return fmt.Errorf("failed to remove session token from Redis: %v", err)
	}
	return nil
}
