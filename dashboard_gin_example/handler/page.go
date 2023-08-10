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
