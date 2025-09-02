package auth

import (
	"os"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/golang-jwt/jwt/v4"
)

var JwtSecretKey string

func InitJwt() {

	JwtSecretKey = os.Getenv("JWT_KEY")

	if JwtSecretKey == "" {
		panic("Could not get hte Jwt Secret Key")
	}
}

func GenerateJWTToken(userId types.UserId, expiration int64) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     expiration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JwtSecretKey)

	return signedToken, err
}
