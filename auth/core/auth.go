package core

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

func GenerateJWT(id, email, signingKey string) string {
	signingKeyBytes := []byte(signingKey)
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()

	tokenString, err := token.SignedString(signingKeyBytes)
	if err != nil {
		return ""
	}
	return tokenString
}

func Authorize(authHeader string, secretKey string) (jwt.Claims, error) {
	signingKey := []byte(secretKey)

	value := strings.Split(authHeader, " ")
	if len(value) != 2 {
		return nil, fmt.Errorf("invalid auth")
	}

	tokenStr := value[1]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
