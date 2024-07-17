package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	v1 := r.Group("/v1")

	v1.POST("/register", controllers.AddUser)
	v1.POST("/login", controllers.LoginUser)

	v2 := r.Group("/v2")
	v2.Use(middleware.TokenMiddleware())

	v2.POST("/update-password", controllers.UpdatePassword)
	v2.POST("/update-pfp", controllers.UpdateProfilePicture)

	r.GET("/ping", controllers.Ping)
	r.GET("/v/user-pfp", controllers.GetUserProfilePicture)
	r.GET("/v/serve-video", controllers.ServerVideo)
	r.POST("/adduser", controllers.AddUser)
	r.POST("/v/add-picture", controllers.UpdateProfilePicture)
	r.POST("/v/add-video", controllers.SaveVideo)
}
