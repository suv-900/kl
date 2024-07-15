package router

import (
	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/controllers"
)

func SetupRouter(e *gin.Engine) {
	e.GET("/ping", controllers.Ping)
	e.POST("/adduser", controllers.AddUser)
	e.POST("/v/add-picture", controllers.UpdateProfilePicture)
}
