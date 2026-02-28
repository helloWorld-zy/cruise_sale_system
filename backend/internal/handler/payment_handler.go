package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PaymentCallbackService 处理来自支付服务商带签名的回调。
type PaymentCallbackService interface {
	HandleCallback(ctx context.Context, provider string, body []byte, signature string) error
}

// PaymentHandler 处理异步的支付服务商回调。
type PaymentHandler struct{ svc PaymentCallbackService }

// NewPaymentHandler 创建 PaymentHandler 实例。
func NewPaymentHandler(svc PaymentCallbackService) *PaymentHandler { return &PaymentHandler{svc: svc} }

// Callback 处理 POST /api/v1/pay/callback 请求。
//
// 安全模型:
//   - 读取完整的请求体（永远不直接存储未受信任的原始输入）。
//   - 从请求头或查询参数中提取特定服务商的签名。
//   - 将签名验证和业务逻辑委托给服务层。
//
// HTTP 约定: 无论结果如何，支付平台都期望返回 HTTP 200；
// 失败信息通过响应体传达（"FAIL" vs {"code":"SUCCESS"}）。
func (h *PaymentHandler) Callback(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.String(http.StatusOK, "FAIL")
		return
	}

	provider := c.Query("provider")
	if provider == "" {
		provider = "wechat" // 默认服务商
	}

	// 微信支付 v3：签名在 Wechatpay-Signature 请求头中。
	// 支付宝：签名在 "sign" 表单/查询参数中。
	signature := c.GetHeader("Wechatpay-Signature")
	if signature == "" {
		signature = c.PostForm("sign")
	}
	if signature == "" {
		// 如果没有签名则拒绝 —— 永远不要处理未签名的回调。
		c.String(http.StatusOK, "FAIL")
		return
	}

	if err := h.svc.HandleCallback(c.Request.Context(), provider, body, signature); err != nil {
		c.String(http.StatusOK, "FAIL")
		return
	}

	// 微信支付 v3 成功响应格式。
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS"})
}
