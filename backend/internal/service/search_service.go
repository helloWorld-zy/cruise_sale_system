package service

type CabinIndexer interface {
	IndexCabin(doc interface{}) error
}

type SearchService struct{ idx CabinIndexer }

func NewSearchService(idx CabinIndexer) *SearchService { return &SearchService{idx: idx} }

func (s *SearchService) IndexCabin(doc interface{}) error {
	return s.idx.IndexCabin(doc)
}
