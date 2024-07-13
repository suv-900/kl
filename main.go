package main

import (
	"github.com/gin-gonic/gin"

	"github.com/suv-900/kl/common"
	"github.com/suv-900/kl/dao"
	"github.com/suv-900/kl/logging"
)

var log = logging.GetLogger()

func main() {

	if err := common.LoadEnv(); err != nil {
		log.Critical(err)
		return
	}
	log.Info("Created Config.")
	if err := dao.Init(); err != nil {
		log.Error(err)
		return
	}
	log.Info("DB connection successfull")

	log.Info("Chores complete")
	engine := gin.New()
	engine.Run()
}
