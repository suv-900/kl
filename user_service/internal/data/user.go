package data

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var (
	ErrRecordNotFound = errors.New("user not found")
	ErrConflict       = errors.New("user already exists")
	ErrUnknown        = errors.New("unknown error occured")
	ErrInternalServer = errors.New("internal server error")
)

type User struct {
	ID        uint64 `gorm:"primarykey"`
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

func (u *User) IsAnonymousUser() bool {
	return *u == *AnonymousUser
}

// pointer reciever avoids copy of the struct
// pointer reciever methods have write access to struct fields
// value reciever methods have read access to struct fields
type UserModel struct {
	DB *gorm.DB
}

func (u UserModel) AddUser(user *User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.WithContext(ctx).Create(user).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return ErrConflict
		default:
			return ErrInternalServer
		}
	}
	return nil
}

func (u UserModel) GetUser(userid uint64) (*User, error) {
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.WithContext(ctx).First(&user, userid).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) UpdateUser(user *User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.WithContext(ctx).Save(user).Error

	return err
}

func (u UserModel) DeleteUser(user *User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.WithContext(ctx).Delete(user).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (u UserModel) CheckUserExists(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := u.DB.WithContext(ctx).Where(&User{Username: username})

	if r.Error != nil {
		return false, r.Error
	}

	return r.RowsAffected > 0, nil
}

func (u UserModel) GetUserPassword(username string) (string, error) {
	var pass string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.WithContext(ctx).Raw(`SELECT password FROM users WHERE username = ?`, username).Scan(&pass).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return pass, ErrRecordNotFound
		default:
			return pass, err
		}
	}

	return pass, nil
}

func (u UserModel) ChangePassword(userid uint64, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.WithContext(ctx).Raw(`UPDATE 
	users SET password = ? WHERE user_id = ?`, password, userid).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (u UserModel) FindActiveUsers() ([]User, error) {
	var users []User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := u.DB.WithContext(ctx).Where("active = ?", true).Find(&users)

	return users, t.Error
}
func (u UserModel) FindSoftDeletedRecords() ([]User, error) {
	var users []User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := u.DB.WithContext(ctx).Where("is_del = 1").Find(&users)
	return users, t.Error
}

type LoginAttemptsResult struct {
	FailedLoginAttempts uint
	FailedLoginTime     time.Time
}

func (u UserModel) GetLoginAttempts(username string) (*LoginAttemptsResult, error) {
	var result LoginAttemptsResult

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.WithContext(ctx).Raw(`SELECT 
	failed_login_attempts,failed_login_time 
	FROM users WHERE username = ?`, username).Scan(&result).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return &result, ErrRecordNotFound
		default:
			return &result, err
		}
	}

	return &result, nil
}

func (u UserModel) UpdateLoginAttempts(username string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.WithContext(ctx).Raw(`UPDATE users SET
	failed_login_attempts = failed_login_attempts + 1,
	failed_login_time = ? WHERE username = ?`, time.Now(), username).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil

}
func (u UserModel) ResetLoginAttempts(username string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := u.DB.WithContext(ctx).Raw(`UPDATE users SET
	failed_login_attempts = 0 WHERE username = ?`, username).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil

}

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
