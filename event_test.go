package splunk_test

import (
	"os"
	"testing"

	"github.com/virzz/splunk-go"
)

func TestSendEvent(t *testing.T) {
	type TestEventData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	forwardTarget := ""
	if v := os.Getenv("SPLUNK_FORWARD_ENDPOINT"); v != "" {
		forwardTarget = eventCfg.Host
		eventCfg.Host = v
	}
	eventClient := splunk.InitEvent(eventCfg).Debug(true)
	if forwardTarget != "" {
		eventClient.Headers(map[string]string{"X-Forward": forwardTarget})
	}
	err := splunk.RawEvent(&TestEventData{
		Title:   "test",
		Content: "test content",
	})
	if err != nil {
		t.Fatal(err)
	}
}
