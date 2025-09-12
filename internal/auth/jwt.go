package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	RoleUser = "user"
)

type UserClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(userID int64, role string, secret string) (string, error) {
	claims := &UserClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "sociotask",
			Subject:   strconv.Itoa(int(userID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func ParseJWTToken(token string, secret string) (*UserClaims, error) {
	claims := &UserClaims{}

	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := t.Claims.(*UserClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
