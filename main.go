package main

import (
	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/router"
)

func main() {
	engine := gin.New()
	router.SetupRouter(engine)

	engine.Run()
}
