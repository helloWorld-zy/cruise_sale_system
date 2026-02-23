package search

import "github.com/meilisearch/meilisearch-go"

type MeiliIndexer struct{ client meilisearch.ServiceManager }

func NewMeiliIndexer(url, key string) *MeiliIndexer {
	return &MeiliIndexer{client: meilisearch.New(url, meilisearch.WithAPIKey(key))}
}

func (m *MeiliIndexer) IndexCabin(doc interface{}) error {
	task, err := m.client.Index("cabins").AddDocuments([]interface{}{doc}, nil)
	_ = task // ignore task for simple adapter
	return err
}
