package main

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/suv-900/kl/config"
	"github.com/suv-900/kl/utils"
)

var log *logging.Logger

func main() {
	//initiate logger
	log = utils.GetLogger()

	c := config.Config{}
	if err := c.LoadEnv(); err != nil {
		log.Critical(err)
		return
	}
	log.Info("config ", c)
	engine := gin.New()
	engine.Run()
}
