package splunk_test

import (
	"testing"

	"github.com/virzz/splunk-go"
)

func TestSearchesList(t *testing.T) {
	res, err := splunk.Searches.List()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
