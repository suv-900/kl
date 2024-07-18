package data

import "gorm.io/gorm"

type Model struct {
	Users interface {
		AddUser(user *User) error
		GetUser(userid uint) (*User, error)
		UpdateUser(user *User) error
		DeleteUser(user *User) error
	}
	Images interface {
		UpdateProfilePicture(image *Image) error
		GetProfilePicture(userid uint) (*Image, error)
	}
}

func GetModel(db *gorm.DB) *Model {
	return &Model{
		Users:  UserModel{DB: db},
		Images: ImageModel{DB: db},
	}
}
