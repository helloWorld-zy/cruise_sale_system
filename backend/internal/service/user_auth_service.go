package service

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrCodeStoreUnavailable 表示验证码存储组件未就绪。
	ErrCodeStoreUnavailable = errors.New("code store unavailable")
	// ErrPhoneOrCodeRequired 表示手机号或验证码缺失。
	ErrPhoneOrCodeRequired = errors.New("phone and code are required")
	// ErrSMSTooFrequent 表示验证码发送过于频繁。
	ErrSMSTooFrequent = errors.New("sms send too frequent")
)

// CodeStore 定义验证码的持久化与校验能力。
type CodeStore interface {
	Save(phone, code string) error
	Verify(phone, code string) bool
}

// UserAuthPolicy 定义短信验证码认证策略参数。
type UserAuthPolicy struct {
	CodeTTL        time.Duration
	ResendInterval time.Duration
	MaxAttempts    int
	LockDuration   time.Duration
	Now            func() time.Time
}

// UserAuthService 提供短信验证码发送、校验与风控能力。
type UserAuthService struct {
	store CodeStore
	now   func() time.Time

	codeTTL        time.Duration
	resendInterval time.Duration
	maxAttempts    int
	lockDuration   time.Duration

	mu             sync.Mutex
	lastSentAt     map[string]time.Time
	expiresAt      map[string]time.Time
	failedAttempts map[string]int
	lockedUntil    map[string]time.Time
}

// NewUserAuthService 使用默认策略创建用户认证服务。
//
//go:noinline
func NewUserAuthService(store CodeStore) *UserAuthService {
	return NewUserAuthServiceWithPolicy(store, UserAuthPolicy{})
}

// NewUserAuthServiceWithPolicy 使用指定策略创建用户认证服务。
//
//go:noinline
func NewUserAuthServiceWithPolicy(store CodeStore, policy UserAuthPolicy) *UserAuthService {
	if policy.CodeTTL <= 0 {
		policy.CodeTTL = 5 * time.Minute
	}
	if policy.ResendInterval <= 0 {
		policy.ResendInterval = time.Minute
	}
	if policy.MaxAttempts <= 0 {
		policy.MaxAttempts = 5
	}
	if policy.LockDuration <= 0 {
		policy.LockDuration = 30 * time.Minute
	}
	if policy.Now == nil {
		policy.Now = time.Now
	}

	return &UserAuthService{
		store:          store,
		now:            policy.Now,
		codeTTL:        policy.CodeTTL,
		resendInterval: policy.ResendInterval,
		maxAttempts:    policy.MaxAttempts,
		lockDuration:   policy.LockDuration,
		lastSentAt:     make(map[string]time.Time),
		expiresAt:      make(map[string]time.Time),
		failedAttempts: make(map[string]int),
		lockedUntil:    make(map[string]time.Time),
	}
}

// SendSMS 发送验证码并记录发送频率与过期时间。
//
//go:noinline
func (s *UserAuthService) SendSMS(phone, code string) error {
	if s.store == nil {
		return ErrCodeStoreUnavailable
	}
	if phone == "" || code == "" {
		return ErrPhoneOrCodeRequired
	}

	now := s.now()
	s.mu.Lock()
	lastSentAt := s.lastSentAt[phone]
	if !lastSentAt.IsZero() && now.Sub(lastSentAt) < s.resendInterval {
		s.mu.Unlock()
		return ErrSMSTooFrequent
	}
	s.mu.Unlock()

	if err := s.store.Save(phone, code); err != nil {
		return err
	}

	s.mu.Lock()
	s.lastSentAt[phone] = now
	s.expiresAt[phone] = now.Add(s.codeTTL)
	s.failedAttempts[phone] = 0
	delete(s.lockedUntil, phone)
	s.mu.Unlock()

	return nil
}

// VerifySMS 校验验证码并根据失败次数执行锁定策略。
//
//go:noinline
func (s *UserAuthService) VerifySMS(phone, code string) bool {
	if s.store == nil || phone == "" || code == "" {
		return false
	}

	now := s.now()
	s.mu.Lock()
	if until := s.lockedUntil[phone]; !until.IsZero() && now.Before(until) {
		s.mu.Unlock()
		return false
	}
	if expireAt := s.expiresAt[phone]; expireAt.IsZero() || now.After(expireAt) {
		s.mu.Unlock()
		return false
	}
	s.mu.Unlock()

	if s.store.Verify(phone, code) {
		s.mu.Lock()
		s.failedAttempts[phone] = 0
		delete(s.lockedUntil, phone)
		s.mu.Unlock()
		return true
	}

	s.mu.Lock()
	s.failedAttempts[phone]++
	if s.failedAttempts[phone] >= s.maxAttempts {
		s.failedAttempts[phone] = 0
		s.lockedUntil[phone] = now.Add(s.lockDuration)
	}
	s.mu.Unlock()

	return false
}

// WechatLogin 处理微信登录流程并返回用户标识。
//
//go:noinline
func (s *UserAuthService) WechatLogin(openID string) string {
	return openID
}
