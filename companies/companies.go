package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Company struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Companies []Company

func parseRawRecord(record string) Company {
	delimiter := ":"
	r := strings.TrimRight(record, delimiter)
	i := strings.LastIndex(r, delimiter)
	return Company{
		Name: r[0:i],
		ID:   r[i+1:],
	}
}

func consumeAndDigestRequest(body string) Companies {
	reader := strings.NewReader(body)
	raw, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	var companies Companies

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

func GetCompanies(c *gin.Context) {
	requestURL := "https://www.sec.gov/Archives/edgar/cik-lookup-data.txt"
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Fatalf("client: could not create request: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/txt")
	req.Header.Set("User-Agent", "Emily me@kimdcottrell.com")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
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
