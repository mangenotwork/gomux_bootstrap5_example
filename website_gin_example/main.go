package main

import (
	"blog_gin_example/routers"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.DebugMode)
	s := routers.Routers()
	s.Run(":18080")
}
