package router

import (
	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/controllers"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/server-status", controllers.ServerStatus)
	r.POST("/adduser", controllers.AddUser)
}
