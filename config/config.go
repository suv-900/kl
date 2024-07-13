package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/suv-900/kl/utils"
)

// Application Configuraton
type Config struct {
	Host       string
	DBName     string
	DBUsername string
	DBPassword string

	BCryptCost int
}

var defaultBcryptCost = 3

var log = utils.GetLogger()

func (c *Config) LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Couldnt load .env:%s", err)
		return err
	}
	var present bool

	c.DBName, present = os.LookupEnv("dbname")
	if !present {
		log.Error("dbname not found in .env file")
		return errors.New("read .env:unsuccessfull")
	}
	c.DBUsername, present = os.LookupEnv("dbusername")
	if !present {
		log.Error("dbusername not found in .env file")
		return errors.New("read .env:unsuccessfull")
	}
	c.DBPassword, present = os.LookupEnv("dbpassword")
	if !present {
		log.Error("dbpassword not found in .env file")
		return errors.New("read .env:unsuccessfull")
	}
	c.Host, present = os.LookupEnv("host")
	if !present {
		log.Error("host not found in .env file")
		return errors.New("read .env:unsuccessfull")
	}

	bcoststr, present := os.LookupEnv("bcrypt_cost")
	if !present {
		//default value 3
		c.BCryptCost = defaultBcryptCost
	} else {

		bcost, err := strconv.Atoi(bcoststr)
		if err != nil {
			log.Error("Unable to convert bcryptcost(string) to int")
			return err
		}
		c.BCryptCost = bcost
	}

	return nil
}
