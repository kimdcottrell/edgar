package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kimdcottrell/edgar/api/framework"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ---
// LOCAL API SERVER.
// Includes router, database, etc
// Inspired from:
//   - https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
//	 - https://github.com/gowebexamples/goreddit/
// ---

type API interface {
	RunMigrations(*Server)
	SetupRoutes(*Server)
}

type Server struct {
	Router
	Database
}

type Database struct {
	*gorm.DB
}

type Router struct {
	*gin.Engine
}

var (
	API_ENDPOINTS = []API{Company{}}
)

const authPath string = "sudo"

func Start() {
	gin.ForceConsoleColor()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: framework.API_DB_CONNECTION_STRING,
	}), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	server := Server{
		Router{
			gin.Default(),
		},
		Database{
			db,
		},
	}
	server.runMigrations()
	server.setupRoutes()
	server.Router.serve()
}

func (s *Server) runMigrations() {
	for _, a := range API_ENDPOINTS {
		a.RunMigrations(s)
	}
}

// routes are to be held in their own respective files, near their handlers
func (s *Server) setupRoutes() {
	for _, a := range API_ENDPOINTS {
		a.SetupRoutes(s)
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
