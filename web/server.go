package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLFiles(
		"web/login.html",
		"web/dashboard.html",
		"web/settings.html",
		"web/inbounds.html",
	)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{})
	})

	r.GET("/inbounds", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inbounds.html", gin.H{})
	})

	r.GET("/settings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "settings.html", gin.H{})
	})

	return r
}
