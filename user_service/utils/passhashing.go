package utils

import (
	"github.com/suv-900/kl/user_service/common"
	"golang.org/x/crypto/bcrypt"
)

var bcryptCost = common.Config.BCryptCost

func GenerateHashedPassword(password []byte) (string, error) {
	var hashedpass string
	hashedbytes, err := bcrypt.GenerateFromPassword(password, bcryptCost)
	if err != nil {
		log.Error(err)
		return hashedpass, err
	}
	return string(hashedbytes), nil
}

// nil on success and err on fail
func ComparePassword(password []byte, dbPassword []byte) error {
	return bcrypt.CompareHashAndPassword(dbPassword, password)
}
