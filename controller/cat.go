package controller

import (
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
