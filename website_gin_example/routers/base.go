package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"website_gin_example/handler"
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

	Router.LoadHTMLGlob("../views/website/*")

	Router.GET("/", handler.List)
	Router.GET("/case1", handler.Index)
	Router.GET("/case2", handler.Index2)

	return Router
}
