package search

import (
	"testing"
)

func TestMeiliIndexer(t *testing.T) {
	// Simple test without actual meilisearch
	indexer := NewMeiliIndexer("http://localhost:7700", "key")
	if indexer == nil {
		t.Errorf("Expected indexer")
	}
	// IndexCabin will fail because connection fails
	err := indexer.IndexCabin(map[string]interface{}{"id": 1})
	if err == nil {
		t.Errorf("Expected error connecting to meilisearch")
	}
}
