package utils

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromJWTClaims(c *gin.Context) (int, error) {
	// Ambil JWT claims dari konteks
	jwtClaims, exists := c.Get("jwtClaims")
	if !exists {
		return 0, errors.New("JWT claims not found in context")
	}

	// Konversi JWT claims ke map[string]interface{}
	claims, ok := jwtClaims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Failed to convert JWT claims to map[string]interface{}")
	}

	// Ambil nilai userID dari JWT claims
	userIDFloat, exists := claims["sub"].(float64)
	if !exists {
		return 0, errors.New("userID not found in JWT claims")
	}

	// Konversi nilai userID dari float64 ke int
	userID := int(userIDFloat)

	return userID, nil
}

func GetUserIDFromJWT(token string) (int, error) {
	// Parse the JWT token
	var tokenString, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Provide the key for verifying the token's signature
		// (replace with your actual key)
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := tokenString.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Failed to convert JWT claims to map[string]interface{}")
	}

	// Ambil nilai userID dari JWT claims
	userIDFloat, exists := claims["sub"].(float64)
	if !exists {
		return 0, errors.New("userID not found in JWT claims")
	}

	// Konversi nilai userID dari float64 ke int
	userID := int(userIDFloat)

	return userID, nil
}
