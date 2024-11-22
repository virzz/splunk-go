package splunk_test

import (
	"testing"

	"github.com/virzz/splunk-go"
)

func TestSearchJob(t *testing.T) {
	q := "search index=main source=soclark"
	sid, err := splunk.Search.Query(q)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sid)
}

func TestSearchDownload(t *testing.T) {
	// splunk.SetDebug(true)
	q := `search index="security" source="sso_log_device" earliest="11/20/2024:00:00:00" latest="11/21/2024:00:00:00" | rename access_time as generate_at | table id username generate_at create_time update_time auth_channel auth_types device_id device_info device_status ua domain error client_device_username os_type os_version public_ip real_ip remote_ip result`
	sid, err := splunk.Search.Query(q)
	if err != nil {
		t.Fatal(err)
	}
	stats, err := splunk.Search.Status(sid, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !stats.IsDone {
		t.Fatal(stats.ResultCount)
	}
	err = splunk.Search.Download(sid, "export.tmp", splunk.OutputModeCSV)
	if err != nil {
		t.Fatal(err)
	}
}
