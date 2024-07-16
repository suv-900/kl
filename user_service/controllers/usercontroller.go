package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/user_service/dao"
	"github.com/suv-900/kl/user_service/logging"
	"github.com/suv-900/kl/user_service/models"
	"github.com/suv-900/kl/user_service/utils"
)

var videosDir = "/home/core/go/kl/user_service/store/videos"
var imagesDir = "/home/core/go/kl/user_service/store/images"

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

	//check login attempts
	loginAttempts, err := dao.GetLoginAttempts(user.Username)
	if err != nil {
		log.Error("error while getting login attempts details ", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if loginAttempts.FailedLoginAttempts > 5 {
		timeNow := time.Now()
		timeDiff := loginAttempts.FailedLoginTime.Sub(timeNow)

		if timeDiff.Abs().Hours() < 2 {
			c.Status(http.StatusUnauthorized)
			return
		} else {
			//set login attempts to 0
			err = dao.ResetLoginAttempts(user.Username)
			if err != nil {
				log.Error("error while reseting login attempts ", err)
				c.Status(http.StatusInternalServerError)
				return
			}
		}
	}

	dbpassword, err := dao.GetUserPassword(user.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}
	if mismatch := utils.ComparePassword([]byte(user.Password), []byte(dbpassword)); mismatch != nil {
		if err := dao.UpdateLoginAttempts(user.Username); err != nil {
			log.Error("error while updating login attempts ", err)
		}
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Status(http.StatusOK)
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
func UpdatePassword(c *gin.Context) {
	token := c.GetHeader("token")

	tokenExpired, userid, tokenInvalid := utils.ValidateToken(token)

	if tokenExpired {
		c.Status(http.StatusBadRequest)
		return
	}
	if tokenInvalid {
		c.Status(http.StatusUnauthorized)
	}

	var password string
	if err := c.ShouldBindJSON(&password); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	password, err := utils.GenerateHashedPassword([]byte(password))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	err = dao.ChangePassword(userid, password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func UpdateUserProfile(c *gin.Context) {
	//get userid
	//handler failure log it out
	value, exists := c.Get("userid")
	if !exists {
		log.Error("userid not found in context map")
		c.Status(http.StatusInternalServerError)
		return
	}

	if value == nil {
		log.Error("Unknown state:userid is nill")
		c.Status(http.StatusInternalServerError)
		return
	}

	var userProfile models.UserProfile

	if err := c.ShouldBindJSON(&userProfile); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

}

func UpdateProfilePicture(c *gin.Context) {
	f, err := c.FormFile("image")
	if err != nil {
		log.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	dest := imagesDir + "/" + f.Filename
	log.Info(dest)
	if err := c.SaveUploadedFile(f, dest); err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	// image := &models.Image{}

	// image.Name = f.Filename
	// image.Size = f.Size
	// image.Location = dest

	// if err := dao.UpdateProfilePicture(image); err != nil {
	// 	log.Error(err)
	// 	c.Status(http.StatusInternalServerError)
	// 	return
	// }

	log.Info("Image added.")
	c.Status(http.StatusOK)
}

func SaveVideo(c *gin.Context) {
	f, err := c.FormFile("video")
	if err != nil {
		log.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	dest := videosDir + "/" + f.Filename
	log.Info(dest)
	if err := c.SaveUploadedFile(f, dest); err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
func ServerVideo(c *gin.Context) {
	c.Status(http.StatusOK)
	c.File(videosDir + "/dogvideo.mp4")
}
func GetUserProfilePicture(c *gin.Context) {
	image, err := dao.GetUserProfilePicture()
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
	c.File(image.Location)
}

// use protobufs
func GetDeletedUsers(c *gin.Context) {
	deletedUsers, err := dao.FindSoftDeletedRecords()
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, deletedUsers)
}

// probably huge
func GetAllActiveUsers(c *gin.Context) {
	activeUsers, err := dao.FindActiveUsers()
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, activeUsers)
}
