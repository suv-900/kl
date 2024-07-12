package models

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   bool
	IsDel    soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
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
