package routing

import (
	"image_service/controllers"

	"github.com/gin-gonic/gin"
)

func SetuptRouter(r *gin.Engine) {
	r.GET("/pfp/:pfp_name", controllers.GetPFP)
}
