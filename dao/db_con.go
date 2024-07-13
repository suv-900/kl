package dao

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host       int
	DBName     string
	DBUsername string
	DBPassword string

	BCryptCost int
}

var db *gorm.DB

// risky
func (c *Config) getvars() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Couldnt load enviorenment vars:%s", err)
		return err
	}

	host, err := strconv.Atoi(os.Getenv("host"))
	if err != nil {
		log.Fatalf("Unable to convert host(string) to uint")
		return err
	}
	c.Host = host

	bcost, err := strconv.Atoi(os.Getenv("bcryptcost"))
	if err != nil {
		log.Fatalf("Unable to convert bcryptcost(string) to uint")
		return err
	}
	c.BCryptCost = bcost

	c.DBName = os.Getenv("dbname")
	c.DBUsername = os.Getenv("dbusername")
	c.DBPassword = os.Getenv("dbpassword")
	return nil
}

func (c *Config) ConnectDB() error {
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
