package dao

import (
	"errors"
	"fmt"

	"github.com/suv-900/kl/logging"
	"github.com/suv-900/kl/models"
)

var log = logging.GetLogger()

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
