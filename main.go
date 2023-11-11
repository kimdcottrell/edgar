package main

import (
	"github.com/gin-gonic/gin"
	c "github.com/kimdcottrell/edgar/companies" // TODO: lazy load this as 'api'
)

func main() {
	r := gin.Default()
	r.GET("/companies", c.GetCompanies)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
