package controller

import (
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/service"
	"encoding/json"
	"fmt"
	"io"
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
		default:
			if err == io.EOF {
				c.JSON(http.StatusBadRequest, gin.H{
					"errors": "Request body is empty",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}
	}

	fmt.Println(matchRequest)
	match, err := controller.matchService.Create(userID, matchRequest)
	if err != nil {
		if err.Error() == ErrCatNotFound.Error() {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Match cat or User cat not found",
			})
			return
		}

		if err.Error() == "THE USER CAT IS NOT BELONG TO THE USER" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat match created succesfully, with id = %v", match.ID),
	})
}

func (controller *matchController) Approve(c *gin.Context) {
	var matchApproval request.MatchApprovalRequest
	jwtClaims, _ := c.Get("jwtClaims")
	claims, _ := jwtClaims.(jwt.MapClaims)
	userID, _ := claims["sub"].(float64)
	err := c.ShouldBindJSON(&matchApproval)

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
		default:
			if err == io.EOF {
				c.JSON(http.StatusBadRequest, gin.H{
					"errors": "Request body is empty",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}
	}

	match, err := controller.matchService.Approval(userID, matchApproval.MatchID, true)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat match with matchId = %v is approved", match),
	})
}

func (controller *matchController) Reject(c *gin.Context) {
	var matchApproval request.MatchApprovalRequest
	jwtClaims, _ := c.Get("jwtClaims")
	claims, _ := jwtClaims.(jwt.MapClaims)
	userID, _ := claims["sub"].(float64)
	err := c.ShouldBindJSON(&matchApproval)

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
		default:
			if err == io.EOF {
				c.JSON(http.StatusBadRequest, gin.H{
					"errors": "Request body is empty",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}
	}

	match, err := controller.matchService.Approval(userID, matchApproval.MatchID, false)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat match with matchId = %v is rejected", match),
	})
}
