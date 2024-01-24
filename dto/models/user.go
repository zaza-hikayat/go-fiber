package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Fullname     string     `json:"fullname" gorm:"column:fullname"`
	Email        string     `json:"email" gorm:"column:email;uniqueIndex"`
	PasswordHash string     `json:"-" gorm:"column:password_hash"`
	PhoneNumber  string     `json:"phone_number" gorm:"column:phone_number"`
	VerifiedAt   *time.Time `json:"verified_at" gorm:"column:verified_at"`
}
