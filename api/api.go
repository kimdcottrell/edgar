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
	SetupRoutes(*Server)
}

type Persistable interface {
	RunMigrations(*Server)
}

type Server struct {
	Router   *gin.Engine
	Database *gorm.DB
}

var (
	API_ENDPOINTS         = []API{Company{}, Ping{}}
	PERSISTABLE_ENDPOINTS = []Persistable{Company{}}
)

func Start() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: framework.API_DB_CONNECTION_STRING,
	}), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	gin.ForceConsoleColor()
	router := gin.Default()

	server := Server{
		router,
		db,
	}
	server.runMigrations()
	server.setupRoutes()
	server.serve()
}

func (s *Server) runMigrations() {
	for _, p := range PERSISTABLE_ENDPOINTS {
		p.RunMigrations(s)
	}
}

// routes are to be held in their own respective files, near their handlers
func (s *Server) setupRoutes() {
	for _, a := range API_ENDPOINTS {
		a.SetupRoutes(s)
	}
}

func (s *Server) serve() {
	h := &http.Server{
		Addr:           ":8080",
		Handler:        s.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	h.ListenAndServe()
	s.Router.Run()
}
