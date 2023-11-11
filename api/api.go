package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ---
// LOCAL API SERVER.
// Includes router, database, etc
// Inspired from:
//   - https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
//	 - https://github.com/gowebexamples/goreddit/
// ---

type Server struct {
	Router
}

type Router struct {
	*gin.Engine
}

const authPath string = "sudo"

func Start() {
	gin.ForceConsoleColor()
	server := Server{
		Router: Router{
			Engine: gin.Default(),
		},
	}
	server.Router.setupRoutes()
	server.Router.serve()
}

// routes are to be held in their own respective files, near their handlers
func (r *Router) setupRoutes() {
	// Simple group: v1
	v1 := r.Group("/v1")
	{
		AddRoutesForCompanies(v1)
		AddRoutesForPing(v1)
	}
}

func (r *Router) serve() {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
	r.Run()
}
