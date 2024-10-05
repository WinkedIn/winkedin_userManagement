package models

import (
	"time"
)

// User represents the user profile in the database
type User struct {
	ID             string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FirstName      string    `gorm:"size:100"`
	LastName       string    `gorm:"size:100"`
	Email          string    `gorm:"size:255;uniqueIndex"`
	PasswordHash   string    `gorm:"size:255"`
	DateOfBirth    time.Time `gorm:"type:date"`
	Gender         string    `gorm:"size:10"`
	LinkedInID     string    `gorm:"size:255;uniqueIndex"`
	CompanyName    string    `gorm:"size:255"`
	JobTitle       string    `gorm:"size:255"`
	Education      string    `gorm:"size:255"`
	ProfilePicture string    `gorm:"type:text"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type UserLoginRequest struct {
	Email       string `json:"email" binding:"required,email"`
	LinkedInJWT string `json:"linkedin_jwt" binding:"required"`
}

type UseLogoutRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
