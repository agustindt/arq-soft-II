package utils

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Claims mínimos que necesitamos
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// obtener secret desde env
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key-here" // fallback (opcional)
	}
	return []byte(secret)
}

// ValidateJWT verifica firma y retorna claims
func ValidateJWT(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		// debe ser HS256 (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	// claims válidos
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
