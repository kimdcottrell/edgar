package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/companies", GetCompanies)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
