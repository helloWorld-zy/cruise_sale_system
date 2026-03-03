package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
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
	// ErrAlipaySignatureInvalid 表示支付宝回调签名校验失败。
	ErrAlipaySignatureInvalid = errors.New("invalid alipay signature")
	// ErrAlipayUIDMismatch 表示客户端上送 UID 与已验签 UID 不一致。
	ErrAlipayUIDMismatch = errors.New("alipay uid mismatch")
	// ErrAlipayPayloadInvalid 表示支付宝登录参数不完整。
	ErrAlipayPayloadInvalid = errors.New("invalid alipay login payload")
	// ErrBindingConfirmationRequired 表示绑定前缺少身份确认。
	ErrBindingConfirmationRequired = errors.New("binding confirmation required")
	// ErrThirdPartyAlreadyBound 表示第三方账号已绑定到其他用户。
	ErrThirdPartyAlreadyBound = errors.New("third-party account already bound")
	// ErrBindPayloadInvalid 表示绑定参数非法。
	ErrBindPayloadInvalid = errors.New("invalid bind payload")
)

const defaultBindConfirmWindow = 5 * time.Minute

// CodeStore 定义验证码的持久化与校验能力。
type CodeStore interface {
	Save(phone, code string) error
	Verify(phone, code string) bool
}

// UserAuthPolicy 定义短信验证码认证策略参数。
type UserAuthPolicy struct {
	CodeTTL          time.Duration
	ResendInterval   time.Duration
	MaxAttempts      int
	LockDuration     time.Duration
	Now              func() time.Time
	AlipaySignSecret string
}

// UserAuthService 提供短信验证码发送、校验与风控能力。
type UserAuthService struct {
	store CodeStore
	now   func() time.Time

	codeTTL        time.Duration
	resendInterval time.Duration
	maxAttempts    int
	lockDuration   time.Duration

	mu                     sync.Mutex
	lastSentAt             map[string]time.Time
	expiresAt              map[string]time.Time
	failedAttempts         map[string]int
	lockedUntil            map[string]time.Time
	alipaySignSecret       string
	bindConfirmWindow      time.Duration
	bindingAuthorizedUntil map[int64]time.Time
	boundAccounts          map[string]int64
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

	alipaySecret := policy.AlipaySignSecret
	if alipaySecret == "" {
		alipaySecret = "MUST-BE-SET-VIA-CONFIG"
	}

	return &UserAuthService{
		store:                  store,
		now:                    policy.Now,
		codeTTL:                policy.CodeTTL,
		resendInterval:         policy.ResendInterval,
		maxAttempts:            policy.MaxAttempts,
		lockDuration:           policy.LockDuration,
		lastSentAt:             make(map[string]time.Time),
		expiresAt:              make(map[string]time.Time),
		failedAttempts:         make(map[string]int),
		lockedUntil:            make(map[string]time.Time),
		alipaySignSecret:       alipaySecret,
		bindConfirmWindow:      defaultBindConfirmWindow,
		bindingAuthorizedUntil: make(map[int64]time.Time),
		boundAccounts:          make(map[string]int64),
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

// AlipayLogin 校验支付宝回调签名并返回已验签 UID。
//
// clientUID 为客户端上送 UID，providerUID 为支付宝回调中的 UID，signature 为回调签名。
func (s *UserAuthService) AlipayLogin(clientUID, providerUID, signature string) (string, error) {
	providerUID = strings.TrimSpace(providerUID)
	signature = strings.TrimSpace(signature)
	clientUID = strings.TrimSpace(clientUID)
	if providerUID == "" || signature == "" {
		return "", ErrAlipayPayloadInvalid
	}
	if !s.verifyAlipaySignature(providerUID, signature) {
		return "", ErrAlipaySignatureInvalid
	}
	if clientUID != "" && clientUID != providerUID {
		return "", ErrAlipayUIDMismatch
	}
	return providerUID, nil
}

func (s *UserAuthService) verifyAlipaySignature(uid, signature string) bool {
	expected := s.signAlipayUID(uid)
	return hmac.Equal([]byte(expected), []byte(strings.ToLower(signature)))
}

func (s *UserAuthService) signAlipayUID(uid string) string {
	mac := hmac.New(sha256.New, []byte(s.alipaySignSecret))
	_, _ = mac.Write([]byte(uid))
	return hex.EncodeToString(mac.Sum(nil))
}

// BindAccount 绑定第三方账号（微信/支付宝/手机号）。
func (s *UserAuthService) BindAccount(userID int64, provider string, identifier string) error {
	provider = strings.TrimSpace(strings.ToLower(provider))
	identifier = strings.TrimSpace(identifier)
	if userID <= 0 || provider == "" || identifier == "" {
		return ErrBindPayloadInvalid
	}

	now := s.now()
	accountKey := provider + ":" + identifier

	s.mu.Lock()
	defer s.mu.Unlock()
	if until := s.bindingAuthorizedUntil[userID]; until.IsZero() || now.After(until) {
		return ErrBindingConfirmationRequired
	}
	if boundUserID, exists := s.boundAccounts[accountKey]; exists && boundUserID != userID {
		return ErrThirdPartyAlreadyBound
	}
	s.boundAccounts[accountKey] = userID
	delete(s.bindingAuthorizedUntil, userID)
	return nil
}

// AuthorizeBinding 通过短信验证码完成绑定前身份确认，验证通过后开启短时绑定窗口。
// phone 为用户绑定的手机号，code 为短信验证码。
func (s *UserAuthService) AuthorizeBinding(userID int64, phone string, code string) error {
	if userID <= 0 || strings.TrimSpace(phone) == "" || strings.TrimSpace(code) == "" {
		return ErrBindPayloadInvalid
	}
	if !s.VerifySMS(phone, code) {
		return ErrBindingConfirmationRequired
	}
	s.mu.Lock()
	s.bindingAuthorizedUntil[userID] = s.now().Add(s.bindConfirmWindow)
	s.mu.Unlock()
	return nil
}
