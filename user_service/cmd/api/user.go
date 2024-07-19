package api

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/user_service/internal/data"
)

var imagesDir = "/home/core/go/kl/user_service/store/images/"

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
func (app *application) UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	var user data.User

	if err := app.models.Users.UpdateUser(&user); err != nil {
		app.internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) GetUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	var userid uint
	//read params
	image, err := app.models.Images.GetProfilePicture(userid)
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	http.ServeFile(w, r, image.Location)
}

func (app *application) UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		app.fileSizeLimitExceded(w, r)
		return
	}

	file, h, err := r.FormFile("image")
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	user := app.contextGetUser(r)

	image := &data.Image{
		UserID:   user.ID,
		Size:     h.Size,
		Location: imagesDir + h.Filename,
	}

	if err := app.models.Images.UpdateProfilePicture(image); err != nil {
		app.internalServerError(w, r)
		return
	}

	destfile, err := os.Create(imagesDir + h.Filename)
	if err != nil {
		app.internalServerError(w, r)
		return
	}
	defer destfile.Close()

	_, err = io.Copy(destfile, file)
	if err != nil {
		app.internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
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
