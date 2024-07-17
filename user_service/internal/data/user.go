package data

import (
	"errors"
	"time"

	"github.com/suv-900/kl/user_service/internal/utils"
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

var AnonymousUser = &User{}

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
	UserID   uint
	User     User `gorm:"constraint:OnDelete:CASCADE;"`
}

type LoginAttemptsResult struct {
	FailedLoginAttempts uint
	FailedLoginTime     time.Time
}

type UserModel struct {
	db *gorm.DB
}

var (
	ErrConflict       = errors.New("user already exists.")
	ErrUnknown        = errors.New("unknown error occured")
	ErrInternalServer = errors.New("internal server error.")
)

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

func (u *UserModel) AddUser(user *User) error {
	var err error
	user.Password, err = utils.GenerateHashedPassword([]byte(user.Password))
	if err != nil {
		return ErrInternalServer
	}
	t := u.db.Create(user)
	return t.Error
}

// sneaky
func (u *UserModel) CheckUserExists(username string) bool {
	r := u.db.Where(&User{Username: username})
	return r.RowsAffected > 0
}

func (u *UserModel) GetLoginAttempts(username string) (*LoginAttemptsResult, error) {
	var result LoginAttemptsResult
	t := u.db.Raw("SELECT failed_login_attempts,failed_login_time FROM users WHERE username = ?", username).Scan(&result)
	return &result, t.Error
}

func (u *UserModel) UpdateLoginAttempts(username string) error {
	t := u.db.Raw(`UPDATE users SET 
	failed_login_attempts = failed_login_attempts + 1,
	failed_login_time = ? WHERE username = ?`, time.Now(), username)
	return t.Error
}
func (u *UserModel) ResetLoginAttempts(username string) error {
	t := u.db.Raw(`UPDATE users SET 
	failed_login_attempts = 0 WHERE username = ?`, username)
	return t.Error
}
func (u *UserModel) GetUserPassword(username string) (string, error) {
	var pass string
	r := u.db.Where("username = ?", username).Select("password").Find(&pass)
	if r.RowsAffected == 0 {
		return "", errors.New("user not found to retreive password")
	}
	return pass, nil
}

// see what are sql.Rows
func (u *UserModel) GetUser(userid uint) (User, error) {
	var user User
	t := u.db.First(&user, userid)
	return user, t.Error
}
func (u *UserModel) ChangePassword(userid uint, password string) error {
	err := u.db.Save(&User{ID: userid, Password: password}).Error
	return err
}
func (u *UserModel) UpdateUser(user User) error {
	t := u.db.Save(&user)
	return t.Error
}

func (u *UserModel) DeleteUser(userid uint) error {
	t := u.db.Delete(&User{}, userid)
	return t.Error
}

// no create for userprofile
func (u *UserModel) UpdateUserProfile(userProfile UserProfile) error {
	t := u.db.Save(userProfile)
	return t.Error
}
func (u *UserModel) GetUserProfilePicture() (*Image, error) {
	var image Image
	t := u.db.Where("id = ?", 10).Find(&image)
	return &image, t.Error
}
func (u *UserModel) UpdateProfilePicture(image *Image) error {
	t := u.db.Save(image)
	return t.Error
}

func (u *UserModel) FindSoftDeletedRecords() ([]User, error) {
	var users []User
	t := u.db.Where("is_del = 1").Find(&users)
	return users, t.Error
}

func (u *UserModel) FindActiveUsers() ([]User, error) {
	var users []User
	t := u.db.Where("active = ?", true).Find(&users)
	return users, t.Error
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
