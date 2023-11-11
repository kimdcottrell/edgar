package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kimdcottrell/edgar/api/middleware"
)

func AddRoutesForPing(g *gin.RouterGroup) {
	// TODO: add real auth middleware
	authorized := g.Group(authPath, middleware.RequireAuth())
	authorized.GET("/ping", getPing)

	g.GET("/ping", getPing)
}

func getPing(c *gin.Context) {
	c.PureJSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}
