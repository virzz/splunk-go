package splunk_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/virzz/splunk-go"
)

var collection = "custom_collection"

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	splunk.New(&splunk.Config{
		Host:     os.Getenv("SPLUNK_ENDPOINT"),
		Username: os.Getenv("SPLUNK_USERNAME"),
		Password: os.Getenv("SPLUNK_PASSWORD"),
	})
}

type UserBaseline struct {
	// Key string `json:"_key"`
	Username string   `json:"username"`
	Country  string   `json:"country"`
	Address  []string `json:"address"`
}

func TestStorageKV(t *testing.T) {
	items := []UserBaseline{
		{Username: "user1", Country: "US"},
		{Username: "user2", Country: "CN"},
		{Username: "user3", Country: "JP"},
		{Username: "user4", Country: "KR"},
		{Username: "user5", Country: "SG"},
	}
	key, err := splunk.Storage.KV.Insert(collection, items[0])
	if err != nil {
		t.Error(err)
	}
	t.Log("Key:", key)
	keys, err := splunk.Storage.KV.Inserts(collection, items[1:])
	if err != nil {
		t.Error(err)
	}
	t.Log("Keys:", keys)
	err = splunk.Storage.KV.Upsert(collection,
		map[string]string{"username": "test"},
		&UserBaseline{Username: "test", Country: "UK", Address: []string{"address1", "address2"}})
	if err != nil {
		t.Error(err)
	}
	var rsp []UserBaseline
	err = splunk.Storage.KV.GetAll(collection, &rsp)
	if err != nil {
		t.Error(err)
	}
	for _, item := range rsp {
		t.Log(item)
	}
}

func TestStorageKVUpsert(t *testing.T) {
	query := map[string]string{"username": "test2"}
	// splunk.SetDebug(true)
	err := splunk.Storage.KV.Upsert(collection,
		query,
		&UserBaseline{Username: "test2", Country: "UK", Address: []string{"address2", "address3"}})
	if err != nil {
		t.Error(err)
	}
	var items []UserBaseline
	err = splunk.Storage.KV.Query(collection,
		query, &items)
	if err != nil {
		t.Error(err)
	}
	for _, item := range items {
		t.Log(item)
	}
}

func TestStorageKVQuery(t *testing.T) {
	splunk.SetDebug(true)
	var items []UserBaseline
	err := splunk.Storage.KV.Query(collection,
		map[string]string{"username": "test"}, &items)
	if err != nil {
		t.Error(err)
	}
	for _, item := range items {
		t.Log(item)
	}
}
