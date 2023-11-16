package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	fw "github.com/kimdcottrell/edgar/api/framework"
)

const (
	cikLookupData fw.SecUrl = "https://www.sec.gov/Archives/edgar/cik-lookup-data.txt"
)

type Company struct {
	ID   string `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Name string `json:"name" gorm:"unique;not null;type:char(10);default:null"`
}

func (c Company) RunMigrations(s *Server) {
	err := s.Database.Migrator().CreateTable(&Company{})
	if err != nil {
		panic("Failed to migrate database")
	}
}

func (c Company) SetupRoutes(s *Server) {
	v1 := s.Router.Group("/v1")
	{
		v1.GET("/companies", getCompanies)
	}
}

func getCompanies(c *gin.Context) {
	res, err := cikLookupData.NewRequest(http.MethodGet, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	b := string(body)
	companies := consumeAndDigestRequest(b)
	jsondata, err := json.Marshal(companies)
	if err != nil {
		log.Fatalf("Error parsing into json data: %s", err)
	}

	c.PureJSON(http.StatusOK, jsondata)
}

func parseRawRecord(record string) Company {
	delimiter := ":"
	r := strings.TrimRight(record, delimiter)
	i := strings.LastIndex(r, delimiter)
	return Company{
		Name: r[0:i],
		ID:   r[i+1:],
	}
}

func consumeAndDigestRequest(body string) []Company {
	reader := strings.NewReader(body)
	raw, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	var companies []Company

	records := strings.Split(string(raw), "\n")
	for _, record := range records {
		if record == "" {
			break
		}
		company := parseRawRecord(record)
		companies = append(companies, company)
	}

	return companies
}
