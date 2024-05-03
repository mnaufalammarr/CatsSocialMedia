package utils

import (
	"errors"

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
