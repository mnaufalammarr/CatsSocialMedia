package controller

import (
	"CatsSocialMedia/model/dto/request"
	"CatsSocialMedia/service"
	"encoding/json"
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

func (*catController) GetById(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat %v created succesfully", "test"),
	})
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
	user, err := controller.catService.Create(catRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User %v created succesfully", user.Name),
	})
}

// func (*catController) Update(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": fmt.Sprintf("Cat %v created succesfully", "test"),
// 	})
// }

// func (*catController) Delete(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": fmt.Sprintf("Cat %v created succesfully", "test"),
// 	})
// }
