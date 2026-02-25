package service

import (
	"errors"
	"sync"
	"testing"
	"time"
)

type flakyIndexer struct {
	mu        sync.Mutex
	failTimes int
	calls     int
}

func (f *flakyIndexer) IndexCabin(_ interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	if f.calls <= f.failTimes {
		return errors.New("temporary error")
	}
	return nil
}

func TestSearchRetryQueueRetrySuccess(t *testing.T) {
	idx := &flakyIndexer{failTimes: 1}
	q := NewSearchRetryQueue(idx, 3, 10)
	q.Start()
	q.Enqueue(map[string]any{"id": 1})

	time.Sleep(1500 * time.Millisecond)

	idx.mu.Lock()
	calls := idx.calls
	idx.mu.Unlock()

	if calls < 2 {
		t.Fatalf("expected retry call, got %d", calls)
	}
}
