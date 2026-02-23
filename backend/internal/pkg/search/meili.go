package search

import "github.com/meilisearch/meilisearch-go"

type MeiliIndexer struct{ client *meilisearch.Client }

func NewMeiliIndexer(url, key string) *MeiliIndexer {
	return &MeiliIndexer{client: meilisearch.NewClient(meilisearch.ClientConfig{Host: url, APIKey: key})}
}

func (m *MeiliIndexer) IndexCabin(doc interface{}) error {
	_, err := m.client.Index("cabins").AddDocuments([]interface{}{doc})
	return err
}
