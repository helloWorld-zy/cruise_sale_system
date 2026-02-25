package search

import (
	"testing"

	"github.com/meilisearch/meilisearch-go"
)

type dummyIndex struct {
	meilisearch.IndexManager
}

func (d dummyIndex) AddDocuments(documents interface{}, options *meilisearch.DocumentOptions) (*meilisearch.TaskInfo, error) {
	return &meilisearch.TaskInfo{TaskUID: 123}, nil
}

type dummyClient struct {
	meilisearch.ServiceManager
}

func (d dummyClient) Index(uid string) meilisearch.IndexManager {
	return dummyIndex{}
}

func TestIndexCabin(t *testing.T) {
	// Use an invalid URL to immediately force a connection failure in AddDocuments
	indexer := NewMeiliIndexer("http://127.0.0.1:0", "dummy-key")
	err := indexer.IndexCabin(map[string]interface{}{"id": 1, "test": "data"})
	if err == nil {
		t.Fatal("Expected an error due to invalid dial address, got nil")
	}

	// Success path
	indexer2 := &MeiliIndexer{client: dummyClient{}}
	err2 := indexer2.IndexCabin(map[string]interface{}{"id": 1})
	if err2 != nil {
		t.Fatal(err2)
	}
}
