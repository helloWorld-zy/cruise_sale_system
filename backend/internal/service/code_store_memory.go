package service

import "sync"

// InMemoryCodeStore 提供进程内验证码存储，用于本地开发和测试环境。
type InMemoryCodeStore struct {
	mu    sync.RWMutex
	codes map[string]string
}

// NewInMemoryCodeStore 创建内存验证码存储。
func NewInMemoryCodeStore() *InMemoryCodeStore {
	return &InMemoryCodeStore{codes: make(map[string]string)}
}

// Save 保存手机号对应的验证码。
func (s *InMemoryCodeStore) Save(phone, code string) error {
	s.mu.Lock()
	s.codes[phone] = code
	s.mu.Unlock()
	return nil
}

// Verify 校验手机号验证码是否匹配。
func (s *InMemoryCodeStore) Verify(phone, code string) bool {
	s.mu.RLock()
	stored := s.codes[phone]
	s.mu.RUnlock()
	return stored != "" && stored == code
}
