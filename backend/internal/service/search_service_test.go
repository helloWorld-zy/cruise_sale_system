package service

import "testing"

type fakeIndexer struct{ called bool }

func (f *fakeIndexer) IndexCabin(_ interface{}) error { f.called = true; return nil }

func TestSearchServiceIndexCabin(t *testing.T) {
	idx := &fakeIndexer{}
	svc := NewSearchService(idx)
	_ = svc.IndexCabin(map[string]interface{}{"id": 1})
	if !idx.called {
		t.Fatal("expected index call")
	}
}
