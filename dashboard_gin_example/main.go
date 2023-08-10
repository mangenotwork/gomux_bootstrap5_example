package main

import (
	"dashboard_gin_example/routers"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.DebugMode)
	s := routers.Routers()
	s.Run(":18082")
}
