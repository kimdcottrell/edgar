package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/kimdcottrell/edgar/api/framework"
	"github.com/kimdcottrell/edgar/api/middleware"
)

type Ping struct{}

func (p Ping) SetupRoutes(s *Server) {

	v1 := s.Router.Group("/v1")
	{
		// TODO: add real auth middleware
		authorized := v1.Group(config.ADMIN_ONLY_PATH, middleware.RequireAuth())
		authorized.GET("/ping", getPing)

		v1.GET("/ping", getPing)
	}

}

func getPing(c *gin.Context) {
	c.PureJSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}
