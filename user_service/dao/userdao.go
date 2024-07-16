package dao

import (
	"errors"
	"fmt"
	"time"

	"github.com/suv-900/kl/user_service/logging"
	"github.com/suv-900/kl/user_service/models"
)

var log = logging.GetLogger()

type LoginAttemptsResult struct {
	FailedLoginAttempts uint
	FailedLoginTime     time.Time
}

// should have updated userid no need to return id
func AddUser(user models.User) error {
	t := db.Create(&user)
	if t.Error != nil {
		log.Fatalf("Couldnt create user:%s", t.Error)
		return t.Error
	}
	fmt.Printf("Rows affected:%d", t.RowsAffected)
	return nil
}

// sneaky
func CheckUserExists(username string) bool {
	r := db.Where(&models.User{Username: username})
	return r.RowsAffected > 0
}
func UpdateUserProfile(userProfile models.UserProfile) error {
	t := db.Save(&userProfile)
	return t.Error
}
func GetLoginAttempts(username string) (*LoginAttemptsResult, error) {
	var result LoginAttemptsResult
	t := db.Raw("SELECT failed_login_attempts,failed_login_time FROM users WHERE username = ?", username).Scan(&result)
	return &result, t.Error
}

func UpdateLoginAttempts(username string) error {
	t := db.Raw(`UPDATE users SET 
	failed_login_attempts = failed_login_attempts + 1,
	failed_login_time = ? WHERE username = ?`, time.Now(), username)
	return t.Error
}
func ResetLoginAttempts(username string) error {
	t := db.Raw(`UPDATE users SET 
	failed_login_attempts = 0 WHERE username = ?`, username)
	return t.Error
}
func GetUserPassword(username string) (string, error) {
	var pass string
	r := db.Where("username = ?", username).Select("password").Find(&pass)
	if r.RowsAffected == 0 {
		log.Error("User not found to retrive password")
		return "", errors.New("user not found to retreive password")
	}
	return pass, nil
}

// see what are sql.Rows
func GetUser(userid uint) (models.User, error) {
	var user models.User
	t := db.First(&user, userid)
	return user, t.Error
}
func ChangePassword(userid uint, password string) error {
	err := db.Save(&models.User{ID: userid, Password: password}).Error
	return err
}
func UpdateUser(user models.User) error {
	t := db.Save(user)
	return t.Error
}

func DeleteUser(userid uint) error {
	t := db.Delete(&models.User{}, userid)
	return t.Error
}

// no create for userprofile
func UpdateUserProfile(userProfile models.UserProfile) error {
	t := db.Save(userProfile)
	return t.Error
}
func GetUserProfilePicture() (*models.Image, error) {
	var image models.Image
	t := db.Where("id = ?", 10).Find(&image)
	return &image, t.Error
}
func UpdateProfilePicture(image *models.Image) error {
	t := db.Save(image)
	return t.Error
}

func FindSoftDeletedRecords() ([]models.User, error) {
	var users []models.User
	t := db.Where("is_del = 1").Find(&users)
	return users, t.Error
}

func FindActiveUsers() ([]models.User, error) {
	var users []models.User
	t := db.Where("active = ?", true).Find(&users)
	return users, t.Error
}
