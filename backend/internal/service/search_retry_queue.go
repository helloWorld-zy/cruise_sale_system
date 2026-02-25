package service

import (
	"sync"
	"time"
)

// retryTask 表示一次待重试的索引任务。
type retryTask struct {
	doc      interface{}
	attempts int
}

// SearchRetryQueue 提供异步索引失败重试能力。
type SearchRetryQueue struct {
	idx       CabinIndexer
	maxRetry  int
	buffer    chan retryTask
	started   bool
	startedMu sync.Mutex
}

// NewSearchRetryQueue 创建索引重试队列并设置默认参数。
func NewSearchRetryQueue(idx CabinIndexer, maxRetry int, queueSize int) *SearchRetryQueue {
	if maxRetry <= 0 {
		maxRetry = 3
	}
	if queueSize <= 0 {
		queueSize = 128
	}
	return &SearchRetryQueue{idx: idx, maxRetry: maxRetry, buffer: make(chan retryTask, queueSize)}
}

// Start 启动后台消费协程；重复调用只会生效一次。
func (q *SearchRetryQueue) Start() {
	q.startedMu.Lock()
	if q.started {
		q.startedMu.Unlock()
		return
	}
	q.started = true
	q.startedMu.Unlock()

	go func() {
		for task := range q.buffer {
			if q.idx == nil {
				continue
			}
			if err := q.idx.IndexCabin(task.doc); err != nil {
				task.attempts++
				if task.attempts < q.maxRetry {
					time.AfterFunc(time.Duration(task.attempts)*time.Second, func() {
						q.Enqueue(task.doc)
					})
				}
			}
		}
	}()
}

// Enqueue 投递待索引文档；队列满时丢弃以避免阻塞主流程。
func (q *SearchRetryQueue) Enqueue(doc interface{}) {
	if q == nil {
		return
	}
	select {
	case q.buffer <- retryTask{doc: doc}:
	default:
	}
}
