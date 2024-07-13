package dao

import (
	"fmt"
	"log"

	"github.com/suv-900/kl/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(c *config.Config) error {
	var err error
	dsn := fmt.Sprintf("host=%d user=%s password=%s dbname=%s", c.Host, c.DBUsername, c.DBPassword, c.DBName)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Couldnt connect to DB:%s", err)
		return err
	}
	return nil
}

// func ValidateRequiredArgs(c *config.Config)error{
// 	switch{
// 	case len(c.DBName) == 0:
// 		errors.New("cannot initiate DB DBName is undefined in Config")
// 	case len(c.DBUsername) == 0:
// 		errors.New("cannot initiate DB DBUsername is undefined in Config")
// 	case len(c.DBPassword) == 0:
// 		errors.New("cannot initiate DB is undefined in Config")
// 	case len(c.DBName) == 0:
// 		errors.New("cannot initiate DB DBName is undefined in Config")
// 	}
// }
