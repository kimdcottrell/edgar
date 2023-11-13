package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	Database
}

type Database struct {
	*sqlx.DB
}

type Router struct {
	*gin.Engine
}

const authPath string = "sudo"

func Start() {
	gin.ForceConsoleColor()
	server := Server{
		Router{
			gin.Default(),
		},
		Database{
			sqlx.MustConnect("postgres", "user=api dbname=api host=persistent_db port=5432 password=password sslmode=disable"),
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
