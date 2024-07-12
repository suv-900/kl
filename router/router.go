package router

import (
	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/controllers"
)

func SetupRouter(e *gin.Engine) {

	e.GET("/server-status", controllers.ServerStatus)
	e.POST("/adduser", controllers.AddUser)
}
