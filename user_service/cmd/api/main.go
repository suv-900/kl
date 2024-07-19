package api

import (
	"github.com/suv-900/kl/user_service/internal/data"
)

type application struct {
	models data.Models
}

func main() {

	if err := LoadEnvVars(); err != nil {
		log.Critical("couldnt read env vars: ", err)
		return
	}

	if err := DBInit(); err != nil {
		log.Error("couldnt initialise database: ", err)
		return
	}
	log.Info("DB connection successfull")

	log.Info("Chores complete")
	log.Info("Creating Gin Engine")

	engine := gin.New()
	SetupRouter(engine)
	engine.Run()
}
