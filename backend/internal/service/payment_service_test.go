package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- 测试替身 ---

type stubGateway struct {
	tradeNo string
	payURL  string
	err     error
}

func (g *stubGateway) CreatePay(_ int64, _ int64) (string, string, error) {
	return g.tradeNo, g.payURL, g.err
}

type stubPayRepo struct {
	payments  map[string]*domain.Payment // 键为 TradeNo
	byID      map[int64]*domain.Payment
	statuses  map[int64]string
	created   []*domain.Payment
	nextID    int64
	createErr error
	findErr   error
}

func newStubPayRepo() *stubPayRepo {
	return &stubPayRepo{
		payments: make(map[string]*domain.Payment),
		byID:     make(map[int64]*domain.Payment),
		statuses: make(map[int64]string),
	}
}

func (r *stubPayRepo) Create(_ context.Context, p *domain.Payment) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.nextID++
	p.ID = r.nextID
	r.created = append(r.created, p)
	r.payments[p.TradeNo] = p
	r.byID[p.ID] = p
	return nil
}

func (r *stubPayRepo) FindByTradeNo(_ context.Context, tradeNo string) (*domain.Payment, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	p, ok := r.payments[tradeNo]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (r *stubPayRepo) FindByID(_ context.Context, id int64) (*domain.Payment, error) {
	p, ok := r.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (r *stubPayRepo) UpdateStatus(_ context.Context, id int64, status string) error {
	r.statuses[id] = status
	if p, ok := r.byID[id]; ok {
		p.Status = status
	}
	return nil
}

type stubBookingStatusRepo struct {
	statuses map[int64]string
}

func newStubBookingStatusRepo() *stubBookingStatusRepo {
	return &stubBookingStatusRepo{statuses: make(map[int64]string)}
}

func (r *stubBookingStatusRepo) UpdateStatus(_ context.Context, id int64, status string) error {
	r.statuses[id] = status
	return nil
}

// fakeVerifier 用于精确控制回调验签与交易号提取分支。
type fakeVerifier struct {
	verifyErr    error
	tradeNo      string
	extractTrade error
}

func (v *fakeVerifier) Verify(_ []byte, _ string) error {
	return v.verifyErr
}

func (v *fakeVerifier) ExtractTradeNo(_ []byte) (string, error) {
	if v.extractTrade != nil {
		return "", v.extractTrade
	}
	return v.tradeNo, nil
}

// --- PaymentService.Create 测试 ---

func TestPaymentServiceCreate_OK(t *testing.T) {
	gw := &stubGateway{tradeNo: "TX001", payURL: "wechat://pay?q=TX001"}
	repo := newStubPayRepo()
	svc := NewPaymentService(gw, repo)

	url, err := svc.Create(context.Background(), 1, 9900, "wechat")

	require.NoError(t, err)
	assert.Equal(t, "wechat://pay?q=TX001", url)
	assert.Len(t, repo.created, 1)
	assert.Equal(t, "TX001", repo.created[0].TradeNo)
	assert.Equal(t, PaymentStatusPending, repo.created[0].Status)
}

func TestPaymentServiceCreate_GatewayError(t *testing.T) {
	gw := &stubGateway{err: errors.New("gateway unavailable")}
	svc := NewPaymentService(gw, newStubPayRepo())

	_, err := svc.Create(context.Background(), 1, 100, "wechat")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gateway.CreatePay")
}

func TestPaymentServiceCreate_RepoCreateError(t *testing.T) {
	gw := &stubGateway{tradeNo: "TX002", payURL: "wechat://pay?q=TX002"}
	repo := newStubPayRepo()
	repo.createErr = errors.New("db down")
	svc := NewPaymentService(gw, repo)

	_, err := svc.Create(context.Background(), 2, 1200, "wechat")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "persist payment")
}

// --- HMACVerifier 测试 ---

func makeHMACSig(t *testing.T, secret string, body []byte) string {
	t.Helper()
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestHMACVerifier_ValidSignature(t *testing.T) {
	v := NewHMACVerifier("mysecret")
	body := []byte(`{"trade_no":"TX001"}`)
	sig := makeHMACSig(t, "mysecret", body)
	assert.NoError(t, v.Verify(body, sig))
}

func TestHMACVerifier_InvalidSignature(t *testing.T) {
	v := NewHMACVerifier("mysecret")
	body := []byte(`{"trade_no":"TX001"}`)
	assert.Error(t, v.Verify(body, "badsig"))
}

func TestHMACVerifier_EmptySignature(t *testing.T) {
	v := NewHMACVerifier("mysecret")
	assert.Error(t, v.Verify([]byte(`{}`), ""))
}

func TestHMACVerifier_ExtractTradeNo_OK(t *testing.T) {
	v := NewHMACVerifier("s")
	body, _ := json.Marshal(map[string]string{"trade_no": "TX999"})
	tradeNo, err := v.ExtractTradeNo(body)
	assert.NoError(t, err)
	assert.Equal(t, "TX999", tradeNo)
}

func TestHMACVerifier_ExtractTradeNo_Missing(t *testing.T) {
	v := NewHMACVerifier("s")
	_, err := v.ExtractTradeNo([]byte(`{}`))
	assert.Error(t, err)
}

func TestHMACVerifier_ExtractTradeNo_InvalidJSON(t *testing.T) {
	v := NewHMACVerifier("s")
	_, err := v.ExtractTradeNo([]byte("not json"))
	assert.Error(t, err)
}

// --- PaymentCallbackServiceImpl 测试 ---

func makeCallbackSvc(payRepo *stubPayRepo, bookRepo *stubBookingStatusRepo, secret string) *PaymentCallbackServiceImpl {
	v := NewHMACVerifier(secret)
	return NewPaymentCallbackService(
		payRepo,
		bookRepo,
		map[string]PaymentVerifier{"wechat": v, "alipay": v},
	)
}

func TestHandleCallback_OK(t *testing.T) {
	secret := "cb-secret"
	payRepo := newStubPayRepo()
	bookRepo := newStubBookingStatusRepo()

	p := &domain.Payment{ID: 10, OrderID: 42, TradeNo: "TX001", Status: PaymentStatusPending, AmountCents: 9900}
	payRepo.payments["TX001"] = p
	payRepo.byID[10] = p

	svc := makeCallbackSvc(payRepo, bookRepo, secret)
	body, _ := json.Marshal(map[string]string{"trade_no": "TX001"})
	sig := makeHMACSig(t, secret, body)

	err := svc.HandleCallback(context.Background(), "wechat", body, sig)

	require.NoError(t, err)
	assert.Equal(t, PaymentStatusPaid, payRepo.statuses[10])
	assert.Equal(t, "paid", bookRepo.statuses[42])
}

func TestHandleCallback_Idempotent_AlreadyPaid(t *testing.T) {
	secret := "cb-secret"
	payRepo := newStubPayRepo()
	bookRepo := newStubBookingStatusRepo()

	p := &domain.Payment{ID: 10, OrderID: 42, TradeNo: "TX001", Status: PaymentStatusPaid}
	payRepo.payments["TX001"] = p
	payRepo.byID[10] = p

	svc := makeCallbackSvc(payRepo, bookRepo, secret)
	body, _ := json.Marshal(map[string]string{"trade_no": "TX001"})
	sig := makeHMACSig(t, secret, body)

	err := svc.HandleCallback(context.Background(), "wechat", body, sig)

	require.NoError(t, err)
	// 已支付：无副作用。
	assert.Empty(t, payRepo.statuses)
	assert.Empty(t, bookRepo.statuses)
}

func TestHandleCallback_InvalidSignature(t *testing.T) {
	svc := makeCallbackSvc(newStubPayRepo(), newStubBookingStatusRepo(), "secret")
	body := []byte(`{"trade_no":"TX001"}`)

	err := svc.HandleCallback(context.Background(), "wechat", body, "badsig")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "verification failed")
}

func TestHandleCallback_UnknownProvider(t *testing.T) {
	svc := makeCallbackSvc(newStubPayRepo(), newStubBookingStatusRepo(), "secret")
	err := svc.HandleCallback(context.Background(), "paypal", []byte(`{}`), "sig")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown payment provider")
}

func TestHandleCallback_TradeNoNotFound(t *testing.T) {
	secret := "secret"
	svc := makeCallbackSvc(newStubPayRepo(), newStubBookingStatusRepo(), secret)
	body, _ := json.Marshal(map[string]string{"trade_no": "NONEXISTENT"})
	sig := makeHMACSig(t, secret, body)

	err := svc.HandleCallback(context.Background(), "wechat", body, sig)
	assert.Error(t, err)
}

func TestHandleCallback_MissingTradeNo(t *testing.T) {
	secret := "secret"
	svc := makeCallbackSvc(newStubPayRepo(), newStubBookingStatusRepo(), secret)
	body := []byte(`{}`)
	sig := makeHMACSig(t, secret, body)

	err := svc.HandleCallback(context.Background(), "wechat", body, sig)
	assert.Error(t, err)
}

func TestHandleCallback_ExtractTradeNoError(t *testing.T) {
	payRepo := newStubPayRepo()
	bookRepo := newStubBookingStatusRepo()
	svc := NewPaymentCallbackService(payRepo, bookRepo, map[string]PaymentVerifier{
		"wechat": &fakeVerifier{tradeNo: "", extractTrade: errors.New("trade extract fail")},
	})

	err := svc.HandleCallback(context.Background(), "wechat", []byte(`{"x":1}`), "sig")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "trade extract fail")
}

func TestHandleCallback_FindByTradeNoError(t *testing.T) {
	payRepo := newStubPayRepo()
	payRepo.findErr = errors.New("query failed")
	bookRepo := newStubBookingStatusRepo()
	svc := NewPaymentCallbackService(payRepo, bookRepo, map[string]PaymentVerifier{
		"wechat": &fakeVerifier{tradeNo: "TX404"},
	})

	err := svc.HandleCallback(context.Background(), "wechat", []byte(`{"trade_no":"TX404"}`), "sig")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "find payment")
}
