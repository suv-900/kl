package main

import (
	"image_service/logging"

	"github.com/gin-gonic/gin"
)

var log = logging.GetLogger()

func main() {
	log.Info("image-service running.")
	engine := gin.New()
	engine.Run()
}
