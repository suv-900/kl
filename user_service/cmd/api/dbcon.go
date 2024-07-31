package api

import (
	"fmt"

	"github.com/suv-900/kl/user_service/internal/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBInit() (*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		Config.Host,
		Config.DBUsername,
		Config.DBPassword,
		Config.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("couldnt connect to DB :%s", err)
		return nil, err
	}
	err = db.AutoMigrate(&data.User{},
		&data.Image{})
	if err != nil {
		log.Errorf("couldnt migrate schemas:%s", err)
		return nil, err
	}

	log.Info("schema migraton successfull.")

	return db, nil
}
