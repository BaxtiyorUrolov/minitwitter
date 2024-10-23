package jwt

import (
	"fmt"
	"time"
	"twitter/config"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(m map[string]interface{}) (string, string, error) {
	// Access token yaratish
	accessToken := jwt.New(jwt.SigningMethodHS256)
	aClaims := accessToken.Claims.(jwt.MapClaims)

	for key, value := range m {
		aClaims[key] = value
	}
	aClaims["exp"] = time.Now().Add(config.AccessExpireTime).Unix() // Access token muddati 20 daqiqa

	// Refresh token yaratish
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rClaims := refreshToken.Claims.(jwt.MapClaims)

	for key, value := range m {
		rClaims[key] = value
	}
	rClaims["exp"] = time.Now().Add(config.RefreshExpireTime).Unix() // Refresh token muddati 24 soat

	// Tokenlarni imzolash
	accessTokenStr, err := accessToken.SignedString(config.SignKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenStr, err := refreshToken.SignedString(config.SignKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

func ExtractClaims(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.SignKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Agar token yaroqli bo'lsa, claimslarni olish
	if ok := token.Valid; ok {
		return token.Claims.(jwt.MapClaims), nil
	}

	return nil, fmt.Errorf("Invalid token")
}
