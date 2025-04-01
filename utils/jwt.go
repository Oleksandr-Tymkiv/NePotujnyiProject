package utils

import (
	"time"

	cf "foodapp/config"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = cf.GetJwtSecret() // უმჯობესია შეინახოთ ENV-ში

func GenerateJWT(id uint, email string) (string, error) {
	// პრეტენზიების შექმნა (მონაცემები, რომლებიც იქნება ჟეტონში)
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // ჟეტონი მოქმედებს 24 საათის განმავლობაში
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
