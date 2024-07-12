package main

import (
	"github.com/gin-gonic/gin"
	"github.com/suv-900/kl/router"
)

func main() {
	r := gin.Default()
	router.SetupRouter(r)

	r.Run()
}
