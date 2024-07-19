package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/user_service/internal/data"
)

var imagesDir = "/home/core/go/kl/user_service/store/images"

func (app *application) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
func (app *application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user data.User
	//implement READJSON
}

func (app *application) LoginUser(w http.ResponseWriter, r *http.Request) {

	var userLogin struct {
		username string `json:"username"`
		password string `json:"password"`
	}

	//READJSON

	//check login attempts
	loginAttempts, err := app.models.Users.GetLoginAttempts(userLogin.username)
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	//avoid hardcoding
	if loginAttempts.FailedLoginAttempts > 5 {
		timeNow := time.Now()
		timeDiff := loginAttempts.FailedLoginTime.Sub(timeNow)

		if timeDiff.Abs().Hours() < 2 {
			//try again after this hours
			return
		} else {
			if err := app.models.Users.ResetLoginAttempts(userLogin.username); err != nil {
				log.Error("error while reseting login attempts ", err)
				app.internalServerError(w, r)
				return
			}
		}
	}

	dbpassword, err := app.models.Users.GetUserPassword(userLogin.username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.userNotFound(w, r)
		default:
			app.internalServerError(w, r)
		}
		return
	}

	if mismatch := app.comparePassword([]byte(userLogin.password), []byte(dbpassword)); mismatch != nil {
		if err := app.models.Users.UpdateLoginAttempts(userLogin.username); err != nil {
			log.Error("error while updating login attempts ", err)
			app.internalServerError(w, r)
			return
		}
		app.wrongcredentials(w, r)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (app *application) CheckUserExists(w http.ResponseWriter, r *http.Request) {
	var username string

	if exists := app.models.Users.CheckUserExists(username); exists {
		app.userExists(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (app *application) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	//READJSON
	var password string

	user := app.contextGetUser(r)

	hashedpassword, err := app.generateHashedPassword([]byte(password))
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	if err := app.models.Users.ChangePassword(user.ID, hashedpassword); err != nil {
		app.internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
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
		log.Error("Unknown state:userid is nil")
		c.Status(http.StatusInternalServerError)
		return
	}

	var userProfile models.UserProfile

	if err := c.ShouldBindJSON(&userProfile); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

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

// get userid from context map
func AddProfilePicture(c *gin.Context) {
	val, found := c.Get("userid")
	if !found {
		log.Error("userid not found in context map")
		return
	}
}

func UpdateProfilePicture(c *gin.Context) {

	var userid uint
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
	image := &models.Image{}

	image.UserID = userid
	image.Name = f.Filename
	image.Size = f.Size
	image.Location = dest

	if err := dao.UpdateProfilePicture(image); err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

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
