package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   bool

	FailedLoginAttempts uint
	FailedLoginTime     time.Time

	IsDel soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}
type UserProfile struct {
	gorm.Model
	Bio       string
	BirthDate time.Time

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

func (u *User) Validate() error {
	switch {
	case len(u.Username) == 0:
		return errors.New("username cannot be empty")
	case len(u.Email) == 0:
		return errors.New("email cannot be empty")
	case len(u.Password) == 0:
		return errors.New("password cannot be empty")
	default:
		return nil
	}
}
