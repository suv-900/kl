package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/config"
	"github.com/suv-900/kl/router"
)

func main() {
	c := config.Config{}
	if err := c.LoadEnv(); err != nil {
		log.Fatal()
		return
	}

	engine := gin.New()
	router.SetupRouter(engine)
	engine.Run()
}
