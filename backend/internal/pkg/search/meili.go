package search

import (
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

// MeiliIndexer 封装了 MeiliSearch 客户端，实现 service.CabinIndexer 接口。
// CRITICAL-02 修复：错误不再被静默忽略，调用方会收到描述性错误信息，
// 可据此决定是否重试或记录故障。
type MeiliIndexer struct{ client meilisearch.ServiceManager }

// NewMeiliIndexer 创建 MeiliSearch 索引器实例。
// url 为 MeiliSearch 服务地址，key 为 API 密钥。
func NewMeiliIndexer(url, key string) *MeiliIndexer {
	return &MeiliIndexer{client: meilisearch.New(url, meilisearch.WithAPIKey(key))}
}

// IndexCabin 将单个舱房文档添加或替换到 MeiliSearch 的 "cabins" 索引中。
// 成功时记录异步任务 ID 以供运维追踪；
// 失败时返回包装后的错误，由调用方处理（记录日志、加入重试队列等）。
func (m *MeiliIndexer) IndexCabin(doc interface{}) error {
	task, err := m.client.Index("cabins").AddDocuments([]interface{}{doc}, nil)
	if err != nil {
		return fmt.Errorf("meilisearch IndexCabin: %w", err)
	}
	// task.TaskUID 可用于状态轮询；如需暴露可通过返回值传递
	_ = task.TaskUID
	return nil
}
