package utils

import (
	"github.com/suv-900/kl/common"
	"golang.org/x/crypto/bcrypt"
)

var bcryptCost = common.Config.BCryptCost

func GenerateHashedPassword(password string) (string, error) {
	var hashedpass string
	hashedbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return hashedpass, err
	}
	return string(hashedbytes), nil
}
