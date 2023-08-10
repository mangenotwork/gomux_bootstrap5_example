package main

import (
	"github.com/gin-gonic/gin"
	"website_gin_example/routers"
)

func main() {

	gin.SetMode(gin.DebugMode)
	s := routers.Routers()
	s.Run(":18081")
}
