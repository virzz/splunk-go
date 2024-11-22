package splunk_test

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/virzz/splunk-go"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	splunk.Init(context.Background(), &splunk.Config{
		Host:     os.Getenv("SPLUNK_ENDPOINT"),
		Username: os.Getenv("SPLUNK_USERNAME"),
		Password: os.Getenv("SPLUNK_PASSWORD"),
	})
}
