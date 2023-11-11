package api

import (
	"io"
	"log"
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
	Router *gin.Engine
}

// routes are to be held in their own respective files, near their handlers
func (s *Server) AddRoutes() {
	s.AddRoutesForCompanies()
}

func (s *Server) Run() {
	s.Router.Run(":8080")
}

// ---
// MIDDLEWARE FOR REQUESTS TO SEC.GOV
// ---

func NewRequest(method string, url string, body io.Reader) (*http.Response, error) {
	// set base timeout to avoid hanging
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalf("Client: could not create request: %s\n", err)
		return req.Response, err
	}

	// read more: https://www.sec.gov/os/webmaster-faq#code-support
	req.Header.Set("Content-Type", "application/txt")
	req.Header.Set("User-Agent", "EDGAR-API me@kimdcottrell.com")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return res, err
	}

	// TODO: pretty sure this will never fire
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d", res.StatusCode)
	}

	return res, err
}
