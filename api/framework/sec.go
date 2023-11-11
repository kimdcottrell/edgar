package framework

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// WRAPPER FOR REQUESTS TO SEC.GOV
type SEC interface {
	NewRequest(string, string, io.Reader) (*http.Response, error)
}

type SecUrl string

func (url SecUrl) validate() (string, error) {
	u := string(url)
	if strings.HasPrefix(u, "https://www.sec.gov") || strings.HasPrefix(u, "https://data.sec.gov") {
		return u, nil
	}
	err := fmt.Errorf("error: %s is not a valid SEC url", u)
	return u, err
}

func (url SecUrl) NewRequest(method string, body io.Reader) (*http.Response, error) {
	// set base timeout to avoid hanging
	u, err := url.validate()
	if err != nil {
		log.Fatalln(err)
		return &http.Response{}, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(method, u, body)
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
