package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/models"
)

func ServerStatus(c *gin.Context) {
}

func AddUser(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		// c.AbortWithStatus(http.StatusUnprocessableEntity)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	err = user.Validate()
	if err != nil {
		// //aborts and stores the error
		// c.AbortWithError(http.StatusBadRequest,err)
		//why even care about such errors
		c.Status(http.StatusBadRequest)
		return
	}

}
