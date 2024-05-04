package controller

import (
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/service"
	"CatsSocialMedia/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userController struct {
	userService service.UserService
}

func NewUserController(service service.UserService) *userController {
	return &userController{service}
}

func (uC *userController) Signup(c *gin.Context) {
	var signupRequest request.SignupRequest
	//var loginRequest request.SignInRequest

	err := c.ShouldBindJSON(&signupRequest)

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

	user, err := uC.userService.Create(signupRequest)

	if err != nil {
		var error string = err.Error()
		if error == "EMAIL ALREADY EXIST" {
			c.JSON(http.StatusConflict, gin.H{
				"errors": "Confilct: email already exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	loginRequest := request.SignInRequest{
		Email:    signupRequest.Email,
		Password: signupRequest.Password,
	}

	tokenString, err := uC.userService.Login(loginRequest)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User Registered Successfully",
		"data": map[string]string{
			"email": user.Email,
			"name":  user.Name,
			"token": tokenString,
		},
	})
}

func (uC *userController) SignIn(c *gin.Context) {
	var loginRequest request.SignInRequest

	err := c.ShouldBindJSON(&loginRequest)

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

	tokenString, err := uC.userService.Login(loginRequest)
	if err != nil {
		if err.Error() == "Invalid email or password" {
			c.JSON(http.StatusNotFound, gin.H{
				"errors": "Invalid email or password",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}

	userId, _ := utils.GetUserIDFromJWT(tokenString)
	user, err := uC.userService.FindByID(userId)

	c.JSON(http.StatusOK, gin.H{
		"message": "User Logged Successfully",
		"data": map[string]string{
			"email": user.Email,
			"name":  user.Name,
			"token": tokenString,
		},
	})
}
