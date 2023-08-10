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
