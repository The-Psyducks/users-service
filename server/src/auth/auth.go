package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecret = "monkeCrack"
	jwtExpirationHours = 2
)

type Claims struct {
	UserId 		string	`json:"user_id"`
	UserAdmin 	bool 	`json:"user_admin"`
	jwt.RegisteredClaims
}

func GenerateToken(userId string, userAdmin bool) (string, error) {
	claims := Claims{
		UserId: userId,
		UserAdmin: userAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * jwtExpirationHours)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "ThePsyducks",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	sign, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error signing JWT token: %w", err)
	}

	return sign, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing JWT token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}
