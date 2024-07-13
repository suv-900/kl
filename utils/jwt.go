package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/suv-900/kl/common"
)

var JWTKEY = common.Config.JWTkey
var tokenExpiryTime = time.Now().Add(60 * time.Minute)

type CustomPayload struct {
	ID uint `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(userid uint) (string, error) {
	payload := CustomPayload{
		ID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := rawToken.SignedString(JWTKEY)
	return token, err
}
