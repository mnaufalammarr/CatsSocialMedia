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

type matchController struct {
	matchService service.MatchService
}

func NewMatchController(service service.MatchService) *matchController {
	return &matchController{service}
}

func (controller *matchController) GetMatches(c *gin.Context) {
	userID, _ := utils.GetUserIDFromJWTClaims(c)
	matches, err := controller.matchService.GetMatches(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    matches,
	})

}

func (controller *matchController) Create(c *gin.Context) {
	var matchRequest request.MatchRequest
	userID, _ := utils.GetUserIDFromJWTClaims(c)
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
		if err != nil {
			// If an error occurred during binding (e.g., no body provided),
			// return a 400 Bad Request status code and an error message
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid input data. Please provide a valid JSON body.",
				"details": err.Error(),
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

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Cat match created succesfully, with id = %v", match.ID),
	})
}

func (controller *matchController) Approve(c *gin.Context) {
	var matchApproval request.MatchApprovalRequest
	userID, _ := utils.GetUserIDFromJWTClaims(c)
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
		}
	}

	match, err := controller.matchService.Approval(userID, matchApproval.MatchID, true)
	if err != nil {
		if err.Error() == "MATCH IS NOT EXIST" {
			c.JSON(http.StatusNotFound, gin.H{
				"errors": "matchId is not found",
			})
			return
		}

		if err.Error() == "MATCHID IS NO LONGER VALID" {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": "matchId is no longer valid",
			})
			return
		}

		if err.Error() == "THE MATCH CAT OWNER IS NOT SAME" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": "unauthorized approval match",
			})
			return
		}

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
	userID, _ := utils.GetUserIDFromJWTClaims(c)
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
		}
	}

	match, err := controller.matchService.Approval(userID, matchApproval.MatchID, false)
	if err != nil {
		if err.Error() == "MATCH IS NOT EXIST" {
			c.JSON(http.StatusNotFound, gin.H{
				"errors": "matchId is not found",
			})
			return
		}

		if err.Error() == "MATCHID IS NO LONGER VALID" {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": "matchId is no longer valid",
			})
			return
		}

		if err.Error() == "THE MATCH CAT OWNER IS NOT SAME" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": "unauthorized reject match",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat match with matchId = %v is rejected", match),
	})
}

func (controller *matchController) Delete(c *gin.Context) {
	matchId := c.Param("id")
	userId, _ := utils.GetUserIDFromJWTClaims(c)

	_, err := controller.matchService.Delete(userId, matchId)
	if err != nil {
		if err.Error() == "MATCH IS NOT EXIST" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Match not found",
			})
			return
		}
		if err.Error() == "UNAUTHORIZED DELETE THIS MATCH" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unathorized delete match",
			})
			return
		}

		if err.Error() == "MATCHID IS ALREADY APPROVED / REJECT" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "matchId is already approved / reject",
			})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Match with matchId %s deleted successfully", matchId),
	})
}
