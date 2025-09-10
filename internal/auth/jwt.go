package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	RoleUser = "user"
)

func GenerateJWTToken(userID int64, role string, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "sociotask",
		"sub":  userID,
		"exp":  time.Now().Add(24 * time.Hour),
		"role": role,
	})
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
