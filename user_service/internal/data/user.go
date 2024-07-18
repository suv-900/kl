package data

import (
	"errors"
	"time"

	"github.com/suv-900/kl/user_service/internal/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var (
	ErrRecordNotFound = errors.New("user not found.")
	ErrConflict       = errors.New("user already exists.")
	ErrUnknown        = errors.New("unknown error occured")
	ErrInternalServer = errors.New("internal server error.")
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Active   bool

	Bio       string
	BirthDate time.Time

	IsDel soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

var AnonymousUser = &User{}

// pointer reciever avoids copy of the struct
// pointer reciever methods have write access to struct fields
// value reciever methods have read access to struct fields
type UserModel struct {
	DB *gorm.DB
}

func (u UserModel) AddUser(user *User) error {
	var err error
	user.Password, err = utils.GenerateHashedPassword([]byte(user.Password))
	if err != nil {
		return ErrInternalServer
	}
	t := u.DB.Create(user)
	return t.Error
}

func (u UserModel) GetUser(userid uint) (*User, error) {
	var user User

	err := u.DB.First(&user, userid).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, ErrInternalServer
		}
	}
	return &user, nil
}

func (u UserModel) UpdateUser(user *User) error {
	t := u.DB.Save(user)
	return t.Error
}

func (u UserModel) DeleteUser(user *User) error {
	t := u.DB.Delete(user)
	return t.Error
}

func (u UserModel) CheckUserExists(username string) bool {
	r := u.DB.Where(&User{Username: username})
	return r.RowsAffected > 0
}

func (u UserModel) GetUserPassword(username string) (string, error) {
	var pass string
	r := u.DB.Where("username = ?", username).Select("password").Find(&pass)
	if r.RowsAffected == 0 {
		return "", errors.New("user not found to retreive password")
	}
	return pass, nil
}

func (u UserModel) ChangePassword(userid uint, password string) error {
	err := u.DB.Save(&User{ID: userid, Password: password}).Error
	return err
}

func (u UserModel) FindActiveUsers() ([]User, error) {
	var users []User
	t := u.DB.Where("active = ?", true).Find(&users)
	return users, t.Error
}
func (u UserModel) FindSoftDeletedRecords() ([]User, error) {
	var users []User
	t := u.DB.Where("is_del = 1").Find(&users)
	return users, t.Error
}

// func pgErrorAnalyser(err error) error {
// 	if err == nil {
// 		return nil
// 	}

// 	if e, ok := err.(*pq.Error); ok {
// 		return pgErrorAnalyser(e.Code)
// 	}
// 	switch {
// 	case err == "23505":
// 		return ErrConflict

// 	default:
// 		return ErrUnknown
// 	}
// }

// sneaky
type LoginAttemptsResult struct {
	FailedLoginAttempts uint
	FailedLoginTime     time.Time
}

func (u UserModel) GetLoginAttempts(username string) (*LoginAttemptsResult, error) {
	var result LoginAttemptsResult
	t := u.DB.Raw("SELECT failed_login_attempts,failed_login_time FROM users WHERE username = ?", username).Scan(&result)
	return &result, t.Error
}

func (u UserModel) UpdateLoginAttempts(username string) error {
	t := u.DB.Raw(`UPDATE users SET
	failed_login_attempts = failed_login_attempts + 1,
	failed_login_time = ? WHERE username = ?`, time.Now(), username)
	return t.Error
}
func (u UserModel) ResetLoginAttempts(username string) error {
	t := u.DB.Raw(`UPDATE users SET
	failed_login_attempts = 0 WHERE username = ?`, username)
	return t.Error
}

// see what are sql.Rows

// no create for userprofile

func (u User) Validate() error {
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
