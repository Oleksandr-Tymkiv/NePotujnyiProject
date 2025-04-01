package utils

import (
	"errors"
	"time"

	cf "foodapp/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

var jwtSecret = cf.GetJwtSecret() // უმჯობესია შეინახოთ ENV-ში

func GenerateJWT(id uint, email string) (string, error) {

	// პრეტენზიების შექმნა (მონაცემები, რომლებიც იქნება ჟეტონში)
	claims := &JWTClaims{
		UserID: id,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "foodapp",
			Subject:   email,
		},
	}

	// ჟეტონის გენერირება
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// ვაწერთ ხელს ჟეტონს
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "ეს სისულელეა, გადააკეთე", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	// ჩატვირთვა კონფიგურაცია
	config, err := cf.LoadConfig()
	if err != nil {
		return nil, err
	}

	// გაანალიზება ნიშანი
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// ხელმოწერის მეთოდის დადასტურება
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// პრეტენზიების დადასტურება და ამოღება
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
