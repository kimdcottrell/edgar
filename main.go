package main

import (
	"github.com/gin-gonic/gin"
	api "github.com/kimdcottrell/edgar/api"
)

func main() {
	gin.ForceConsoleColor()
	server := api.Server{
		Router: gin.Default(),
	}
	server.AddRoutes()
	server.RunRouter()

}
