package utils

import (
	"github.com/walk1ng/gin-photo-gallery-storage/conf"
	"github.com/walk1ng/gin-photo-gallery-storage/constant"
	"go.uber.org/zap"

	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserClaim struct
type UserClaim struct {
	UserName string `json:"userName"`
	jwt.StandardClaims
}

// GenerateJWT func to gen a JWT string based on the user name
func GenerateJWT(userName string) (string, error) {
	// define a user claim
	claim := UserClaim{
		userName,
		jwt.StandardClaims{
			Issuer:    constant.PhotoStorageAdmin,
			ExpiresAt: time.Now().Add(constant.JwtExpMinute * time.Minute).Unix(),
		},
	}

	// generate the claim and the digital signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtString, err := token.SignedString([]byte(conf.ServerCfg.Get(constant.JwtSecret)))
	if err != nil {
		AppLogger.Fatal(err.Error(), zap.String("service", "GenerateJWT()"))
		return "", err
	}
	return jwtString, nil
}

// ParseJWT func to parse a JWT into a user claim
func ParseJWT(jwtString string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(jwtString, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.ServerCfg.Get(constant.JwtSecret)), nil
	})

	if token != nil && err == nil {
		if claim, ok := token.Claims.(*UserClaim); ok && token.Valid {
			return claim, nil
		}
	}
	return nil, err
}
