package splunk

import (
	"encoding/json"
)

type kv struct {
	owner, app string
	collection string
}

const (
	// KVStore is the key-value store
	kvBaseUri   = "/servicesNS/{owner}/{app}/storage/collections/data/{collection}"
	kvMultiUri  = "/servicesNS/{owner}/{app}/storage/collections/data/{collection}/batch_save"
	kvUpdateUri = "/servicesNS/{owner}/{app}/storage/collections/data/{collection}/{_key}"
)

func (s *kv) Set(owner, app string) {
	s.app = app
	s.owner = owner
}

// Insert an item into the {collection}.
func (s *kv) Insert(collection string, data any) (string, error) {
	var key struct {
		Key string `json:"_key"`
	}
	_, err := client.R().
		SetPathParams(map[string]string{
			"owner":      s.owner,
			"app":        s.app,
			"collection": collection,
		}).
		SetDebug(clientDebug).
		SetBody(data).
		SetResult(&key).
		Post(kvBaseUri)
	if err != nil {
		return "", err
	}
	return key.Key, nil
}

// Update the record for {_key} ID in {collection}.
func (s *kv) Update(collection, key string, data any) error {
	_, err := client.R().
		SetPathParams(map[string]string{
			"owner":      s.owner,
			"app":        s.app,
			"collection": collection,
			"_key":       key,
		}).
		SetDebug(clientDebug).
		SetBody(data).
		Post(kvUpdateUri)
	if err != nil {
		return err
	}
	return nil
}

// Inserts insert multiple items into the {collection}.
func (s *kv) Inserts(collection string, data any) ([]string, error) {
	var keys []string
	_, err := client.R().
		SetPathParams(map[string]string{
			"owner":      s.owner,
			"app":        s.app,
			"collection": collection,
		}).
		SetDebug(clientDebug).
		SetBody(data).
		SetResult(&keys).
		Post(kvMultiUri)
	if err != nil {
		return keys, err
	}
	return keys, nil
}

// Query retrieves items from the {collection} based on a query.
func (s *kv) Query(collection string, query map[string]string, rsp any) error {
	req := client.R().
		SetPathParams(map[string]string{
			"owner":      s.owner,
			"app":        s.app,
			"collection": collection,
		}).
		SetDebug(clientDebug)
	if query != nil && len(query) > 0 {
		queryBuf, err := json.Marshal(query)
		if err != nil {
			return err
		}
		req.SetQueryParam("query", string(queryBuf))
	}
	if rsp != nil {
		req.SetResult(&rsp)
	}
	_, err := req.Get(kvBaseUri)
	if err != nil {
		return err
	}
	return nil
}

// GetAll retrieves all items from the {collection}.
func (s *kv) GetAll(collection string, rsp any) error {
	req := client.R().
		SetPathParams(map[string]string{
			"owner":      s.owner,
			"app":        s.app,
			"collection": collection,
		}).
		SetDebug(clientDebug)
	if rsp != nil {
		req.SetResult(&rsp)
	}
	_, err := req.Get(kvBaseUri)
	if err != nil {
		return err
	}
	return nil
}

func (s *kv) Upsert(collection string, query map[string]string, data any) error {
	type kvData struct {
		Key string `json:"_key"`
	}
	var items = []kvData{}
	err := s.Query(collection, query, &items)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		_, err = s.Insert(collection, data)
	} else {
		for _, item := range items {
			err = s.Update(collection, item.Key, data)
			if err != nil {
				return err
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}
