package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type catController struct {
	// catService service.CatService
}

func NewCatController(
// service service.CatService
) *catController {
	return &catController{
		// service
	}
}

func (*catController) All(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Cat %v created succesfully", "test"),
	})
}
