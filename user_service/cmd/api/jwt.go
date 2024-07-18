package api

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// var JWTKEY = common.Config.JWTkey
var jwtKey = "JWTKEY"
var tokenExpiryTime = time.Now().Add(60 * time.Minute)
var ErrTokenInvalid = errors.New("token invalid")

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
	token, err := rawToken.SignedString(jwtKey)
	return token, err
}

// i dont even know what i want
// tokenExpired id tokenInvalid
func VerifyToken(token string) (uint, error) {
	var userid uint
	t, err := jwt.ParseWithClaims(token, &CustomPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		//this probably means token is expired
		return userid, err
	}
	if p, ok := t.Claims.(*CustomPayload); ok && t.Valid {
		//this probably means token is invalid
		userid = p.ID
		return userid, nil
	}
	return userid, ErrTokenInvalid
}
