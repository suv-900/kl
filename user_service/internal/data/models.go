package data

import "gorm.io/gorm"

type Models struct {
	Users interface {
		AddUser(user *User) error
		GetUser(userid uint) (*User, error)
		UpdateUser(user *User) error
		DeleteUser(user *User) error

		GetUserPassword(username string) (string, error)
		CheckUserExists(username string) bool
		ChangePassword(userid uint, password string) error

		GetLoginAttempts(username string) (*LoginAttemptsResult, error)
		ResetLoginAttempts(username string) error
		UpdateLoginAttempts(username string) error

		FindActiveUsers() ([]User, error)
		FindSoftDeletedRecords() ([]User, error)
	}
	Images interface {
		UpdateProfilePicture(image *Image) error
		GetProfilePicture(userid uint) (*Image, error)
	}
}

func GetModel(db *gorm.DB) *Models {
	return &Models{
		Users:  UserModel{DB: db},
		Images: ImageModel{DB: db},
	}
}
