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
	Binary []byte
	Name   string
	Size   int64

	UserID uint
	User   User `gorm:"constraint:OnDelete:CASCADE;"`
}
