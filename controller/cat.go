package controller

import (
	"CatsSocialMedia/model"
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/service"
	"CatsSocialMedia/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type catController struct {
	catService service.CatService
}

func NewCatController(service service.CatService) *catController {
	return &catController{service}
}

func (*catController) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat %v created succesfully", "test"),
	})
}

func (controller *catController) FindAll(c *gin.Context) {
	filterParams := make(map[string]interface{})

	userID, _ := utils.GetUserIDFromJWTClaims(c)

	// Parse query parameters
	for key, values := range c.Request.URL.Query() {
		value := values[0] // We only use the first value if there are multiple values for the same key
		switch key {
		case "id":
			filterParams["id"] = value
		case "limit":
			limit, err := strconv.Atoi(value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'limit'"})
				return
			}
			filterParams["limit"] = limit
		case "offset":
			offset, err := strconv.Atoi(value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'offset'"})
				return
			}
			filterParams["offset"] = offset
		case "race":
			filterParams["race"] = value
		case "sex":
			filterParams["sex"] = value
		case "hasMatched":
			hasMatched, err := strconv.ParseBool(value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'hasMatched'"})
				return
			}
			filterParams["hasMatched"] = hasMatched
		case "ageInMonth":
			filterParams["ageInMonth"] = value
		case "owned":
			owned, err := strconv.ParseBool(value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'owned'"})
				return
			}
			filterParams["owned"] = owned
		case "search":
			filterParams["search"] = value
			// Add parsing for other filters similarly...
		}
	}
	filterParams["userID"] = userID
	fmt.Println(filterParams)
	// Call service to get cats with filters
	cats, err := controller.catService.FindAll(filterParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, cats)
}

func (controller *catController) FindByUserID(c *gin.Context) {
	// Retrieve user ID from request or any other source
	userID := 123 // Example: replace with actual user ID retrieval logic

	cat, err := controller.catService.FindByUserID(userID)
	if err != nil {
		// Handle error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Check if cat is found
	if _, ok := cat.(model.Cat); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func (controller *catController) FindByID(c *gin.Context) {
	catID := c.Param("id")

	// Call service to find cat by ID
	cat, err := controller.catService.FindByID(catID)
	if err != nil {
		if errors.Is(err, errors.New("cat not found")) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cat not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, cat)
}

func (controller *catController) Create(c *gin.Context) {

	var catRequest request.CatRequest
	userID, _ := utils.GetUserIDFromJWTClaims(c)
	err := c.ShouldBindJSON(&catRequest)

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
	fmt.Println(userID)
	catRequest.UserId = int(userID)
	fmt.Println(catRequest)
	cat, err := controller.catService.Create(catRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat %v created succesfully", cat.Name),
	})
}

func (controller *catController) Update(c *gin.Context) {
	catID := c.Param("id")

	// Bind request body to CatRequest struct
	var catRequest request.CatRequest
	if err := c.ShouldBindJSON(&catRequest); err != nil {
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
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": err.Error(),
			})
			return
		}
	}

	// Call service to update the cat
	cat, err := controller.catService.Update(catID, catRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat %v updated successfully", cat.Name),
	})
}

func (controller *catController) Delete(c *gin.Context) {
	catID := c.Param("id")

	// Call service to delete cat by ID
	err := controller.catService.Delete(catID)
	if err != nil {
		if errors.Is(err, errors.New("cat not found")) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cat not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cat deleted successfully",
	})
}
