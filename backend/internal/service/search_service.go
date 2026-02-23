package service

// CabinIndexer 定义舱房文档索引的端口接口。
// 由搜索引擎适配器（如 MeiliSearch）提供具体实现。
type CabinIndexer interface {
	IndexCabin(doc interface{}) error // 将舱房文档索引到搜索引擎
}

// SearchService 提供搜索引擎索引相关的业务逻辑。
type SearchService struct{ idx CabinIndexer }

// NewSearchService 创建搜索服务实例。
func NewSearchService(idx CabinIndexer) *SearchService { return &SearchService{idx: idx} }

// IndexCabin 将舱房数据索引到搜索引擎，以支持全文搜索功能。
func (s *SearchService) IndexCabin(doc interface{}) error {
	return s.idx.IndexCabin(doc)
}
