package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	UserID     uint
	User       User
	Bio        *string
	ProfilePic *pq.ByteaArray
	BirthDate  *time.Time
}
