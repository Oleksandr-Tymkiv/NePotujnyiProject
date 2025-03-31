package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-super-secret-key") // Краще зберігати в ENV

func GenerateJWT(id uint, email string) (string, error) {
	// Створюємо claims (дані, які будуть в токені)
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Токен діє 24 години
	}

	// Генеруємо токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Підписуємо токен
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
