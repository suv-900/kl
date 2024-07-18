package api

import (
	"fmt"

	"github.com/suv-900/kl/user_service/internal/data"
	"github.com/suv-900/kl/user_service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Models *data.Models

func Init() error {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		common.Config.Host,
		common.Config.DBUsername,
		common.Config.DBPassword,
		common.Config.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("couldnt connect to DB :%s", err)
		return err
	}
	err = db.AutoMigrate(&models.User{},
		&models.Image{},
		&models.UserProfile{})
	if err != nil {
		log.Errorf("couldnt migrate schemas:%s", err)
		return err
	}

	log.Info("schema migraton successfull.")

	setModel(db)

	return nil
}

func setModel(db *gorm.DB) {
	Models = data.GetModel(db)
}
