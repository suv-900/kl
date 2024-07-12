package models

import (
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	Active   bool
	IsDel    soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}
