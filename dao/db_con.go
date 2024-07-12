package dao

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// risky
func getvars() ([4]string, error) {
	var vars [4]string
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Couldnt load enviorenment vars:%s", err)
		return vars, err
	}

	vars[0] = os.Getenv("host")
	vars[1] = os.Getenv("user")
	vars[2] = os.Getenv("password")
	vars[3] = os.Getenv("dbname")
	fmt.Println(vars)

	return vars, nil
}

func ConnectDB() error {
	vars, err := getvars()
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", vars[0], vars[1], vars[2], vars[3])

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Couldnt connect to DB:%s", err)
		return err
	}
	return nil
}
