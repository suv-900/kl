package dao

import (
	"fmt"
	"log"

	"github.com/suv-900/kl/models"
)

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

func FindSoftDeletedRecords() ([]models.User, error) {
	var users []models.User
	t := db.Where("is_del = 1").Find(&users)
	return users, t.Error
}
