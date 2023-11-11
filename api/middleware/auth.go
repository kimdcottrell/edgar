package middleware

import "github.com/gin-gonic/gin"

// TODO: add real auth middleware
func RequireAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"foo": "bar",
	})
}
