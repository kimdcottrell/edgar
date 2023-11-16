package framework

import (
	"fmt"
	"os"
)

var (
	API_DB_CONNECTION_STRING = fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=%s",
		os.Getenv("API_DB_USER"),
		os.Getenv("API_DB_PASSWORD"),
		os.Getenv("API_DB_HOST"),
		os.Getenv("API_DB_NAME"),
		"disable",
	)
)
