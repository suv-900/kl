package models

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	Bio       *string
	BirthDate *time.Time

	UserID uint
	User   User `gorm:"constraint:OnDelete:CASCADE;"`
}

type Image struct {
	gorm.Model
	Name     string
	Size     int64
	Location string
	// UserID uint
	// User   User `gorm:"constraint:OnDelete:CASCADE;"`
}
