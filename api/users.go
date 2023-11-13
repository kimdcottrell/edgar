package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/kimdcottrell/edgar/api/middleware"
)

type User struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Username string    `db:"username" json:"username"`
	Password string    `db:"password" json:"password"`
}

func AddRoutesForUsers(g *gin.RouterGroup) {
	// TODO: add real auth middleware
	authorized := g.Group(authPath, middleware.RequireAuth())
	authorized.GET("/user", getPing)

	g.GET("/user", getPing)
}

func getUsers(c *gin.Context) {
	c.PureJSON(http.StatusOK, gin.H{
		"ping": "pong",
	})
}
