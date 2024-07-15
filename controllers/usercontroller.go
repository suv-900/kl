package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/dao"
	"github.com/suv-900/kl/logging"
	"github.com/suv-900/kl/models"
	"github.com/suv-900/kl/utils"
)

var log = logging.GetLogger()

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
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
	username := c.Param("username")
	if len(username) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}
	if exists := dao.CheckUserExists(username); exists {
		c.Status(http.StatusConflict)
		return
	} else {
		c.Status(http.StatusOK)
		return
	}
}

func UpdateProfilePicture(c *gin.Context) {
	f, err := c.FormFile("image")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	image := &models.Image{}

	image.Name = f.Filename
	image.Size = f.Size

	file, err := f.Open()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	image.Binary = make([]byte, image.Size)
	_, err = file.Read(image.Binary)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := dao.UpdateProfilePicture(image); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Info("Image added.")
	c.Status(http.StatusOK)
}

func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	if len(user.Username) == 0 || len(user.Password) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	dbpassword, err := dao.GetUserPassword(user.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}
	if mismatch := utils.ComparePassword([]byte(user.Password), []byte(dbpassword)); mismatch != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Status(http.StatusOK)
}

// c.SetCookie("token", s, 3600, "/", "localhost", false, true)
