package routers

import (
	"dashboard_gin_example/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
}

func Routers() *gin.Engine {
	Router.StaticFS("/static", http.Dir("../static"))

	//模板
	// 自定义模板方法
	//Router.SetFuncMap(template.FuncMap{
	//	"formatAsDate": FormatAsDate,
	//})

	Router.Delims("{[", "]}")

	Router.LoadHTMLGlob("../views/dashboard/*")

	Router.GET("/", handler.List)
	Router.GET("/case1", handler.Index)
	Router.GET("/case2", handler.Index2)
	Router.GET("/case3", handler.Index3)
	Router.GET("/case4", handler.Index4)
	Router.GET("/case5", handler.Index5)
	Router.GET("/case6", handler.Index6)
	Router.GET("/case7", handler.Index7)
	Router.GET("/case8", handler.Index8)
	Router.GET("/case9", handler.Index9)
	Router.GET("/case10", handler.Index10)

	return Router
}
