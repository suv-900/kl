package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/suv-900/kl/user_service/common"
	"github.com/suv-900/kl/user_service/logging"
)

var JWTKEY = common.Config.JWTkey
var tokenExpiryTime = time.Now().Add(60 * time.Minute)

var log = logging.GetLogger()

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

// i dont even know what i want
// tokenExpired id tokenInvalid
func ValidateToken(token string) (bool, uint, bool) {
	var id uint
	t, err := jwt.ParseWithClaims(token, &CustomPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTKEY), nil
	})
	if err != nil {
		log.Error(err)
		return true, id, false
	}
	if p, ok := t.Claims.(*CustomPayload); ok && t.Valid {
		id = p.ID
		return false, id, false
	}
	return false, id, true
}
