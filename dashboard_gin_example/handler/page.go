package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func List(c *gin.Context) {
	c.HTML(http.StatusOK, "list.html", gin.H{})
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func Index2(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case2.html", gin.H{})
}

func Index3(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case3.html", gin.H{})
}

func Index4(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case4.html", gin.H{})
}

func Index5(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case5.html", gin.H{})
}

func Index6(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case6.html", gin.H{})
}

func Index7(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case7.html", gin.H{})
}

func Index8(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case8.html", gin.H{})
}

func Index9(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case9.html", gin.H{})
}

func Index10(c *gin.Context) {
	c.HTML(http.StatusOK, "index_case10.html", gin.H{})
}
