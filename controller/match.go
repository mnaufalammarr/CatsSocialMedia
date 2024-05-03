package controller

import (
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type matchController struct {
	matchService service.MatchService
}

func NewMatchController(service service.MatchService) *matchController {
	return &matchController{service}
}

func (controller *matchController) Create(c *gin.Context) {
	var matchRequest request.MatchRequest
	jwtClaims, _ := c.Get("jwtClaims")
	claims, _ := jwtClaims.(jwt.MapClaims)
	userID, _ := claims["sub"].(float64)
	err := c.ShouldBindJSON(&matchRequest)

	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errorMessages := []string{}
			for _, e := range err.(validator.ValidationErrors) {
				errorMessage := fmt.Sprintf("Error on field: %s, condition: %s", e.Field(), e.ActualTag())
				errorMessages = append(errorMessages, errorMessage)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errorMessages,
			})
			return
		case *json.UnmarshalTypeError:
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}
	}

	fmt.Println("request payload")
	fmt.Println(matchRequest)
	match, err := controller.matchService.Create(userID, matchRequest)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat match created succesfully, with id = %v", match.ID),
	})
}