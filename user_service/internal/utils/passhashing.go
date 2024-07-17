package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPassword(password []byte) (string, error) {
	var hashedpass string
	hashedbytes, err := bcrypt.GenerateFromPassword(password, 3)
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
