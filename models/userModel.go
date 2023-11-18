package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username         string    `gorm:"unique"`
	Email            string    `gorm:"unique"`
	PasswordHash     string
	Token            string
	RefreshToken     string
	RegistrationTime time.Time
	LastLoginTime    time.Time
	ProfilePicLink   string
}

type Registration struct {
	Username string
	Email    string
	Password string
}

type Login struct {
	UsernameOrEmail string
	Password        string
}
