package api

import (
	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/user_service/logging"
)

var log = logging.GetLogger()

func main() {
	if err := dao.Init(); err != nil {
		log.Error(err)
		return
	}
	log.Info("DB connection successfull")

	log.Info("Chores complete")
	log.Info("Creating Gin Engine")
	engine := gin.New()
	router.SetupRouter(engine)
	engine.Run()
}
