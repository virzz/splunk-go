package splunk_test

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/virzz/splunk-go"
)

var eventCfg *splunk.EventConfig

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

	eventCfg = &splunk.EventConfig{
		Host:   os.Getenv("SPLUNK_HOST"),
		Index:  os.Getenv("SPLUNK_INDEX"),
		Source: os.Getenv("SPLUNK_SOURCE"),
		Token:  os.Getenv("SPLUNK_TOKEN"),
	}
}
