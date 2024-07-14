package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/dao"
	"github.com/suv-900/kl/models"
	"github.com/suv-900/kl/utils"
)

func ServerStatus(c *gin.Context) {
}

func AddUser(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	err = user.Validate()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user.Password, err = utils.GenerateHashedPassword([]byte(user.Password))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := dao.AddUser(user); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	var token string
	if token, err = utils.GenerateToken(user.ID); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	//fix
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.Status(http.StatusCreated)
}

func CheckUserExists(c *gin.Context) {

}

// c.SetCookie("token", s, 3600, "/", "localhost", false, true)
