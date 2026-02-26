package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/middleware"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// UserAuthService 定义用户登录验证码校验能力。
type UserAuthService interface {
	VerifySMS(phone, code string) bool
	SendSMS(phone, code string) error
}

// UserRepository 提供 C 端用户数据库操作接口。
type UserRepository interface {
	FindOrCreateByPhone(phone string) (*domain.User, error)
}

// UserHandler 处理 C 端登录、发码与个人信息接口。
type UserHandler struct {
	authSvc   UserAuthService
	userRepo  UserRepository // 可为 nil（向下兼容）
	jwtSecret string
}

// NewUserHandler 创建用户处理器实例。
func NewUserHandler(authSvc UserAuthService, jwtSecret string) *UserHandler {
	return &UserHandler{authSvc: authSvc, jwtSecret: jwtSecret}
}

// NewUserHandlerWithRepo 在 NewUserHandler 基础上注入 UserRepository，
// 使 Login 存储 user.ID（整型）而非手机号作为 JWT sub（M-02 修复）。
func NewUserHandlerWithRepo(authSvc UserAuthService, userRepo UserRepository, jwtSecret string) *UserHandler {
	return &UserHandler{authSvc: authSvc, userRepo: userRepo, jwtSecret: jwtSecret}
}

// UserLoginRequest 表示用户短信验证码登录请求体。
type UserLoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// UserLoginResponse 表示登录成功后的令牌返回结构。
type UserLoginResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

// Login 校验短信验证码并签发 JWT 令牌。
func (h *UserHandler) Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	if h.authSvc == nil || h.jwtSecret == "" {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "user auth service unavailable")
		return
	}

	if !h.authSvc.VerifySMS(req.Phone, req.Code) {
		response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "invalid sms code")
		return
	}

	// M-02 修复：JWT sub 存储数据库 user.ID（整型字符串），而非手机号
	sub := req.Phone // 兼容旧模式（无 userRepo 时）
	if h.userRepo != nil {
		user, err := h.userRepo.FindOrCreateByPhone(req.Phone)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "failed to lookup user")
			return
		}
		sub = strconv.FormatInt(user.ID, 10)
	}

	expireAt := time.Now().Add(24 * time.Hour)
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   sub,
		"roles": []string{"user"},
		"exp":   expireAt.Unix(),
	})
	token, _ := tokenObj.SignedString([]byte(h.jwtSecret))

	response.Success(c, UserLoginResponse{Token: token, ExpireAt: expireAt})
}

// SendCodeRequest 表示短信验证码发送请求体。
type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// SendCode 发送短信验证码（POST /users/sms-code）。
func (h *UserHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if h.authSvc == nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "user auth service unavailable")
		return
	}
	// 服务端生成 6 位 OTP，避免客户端自派发验证码（安全要求）
	n, _ := rand.Int(rand.Reader, big.NewInt(1_000_000))
	code := fmt.Sprintf("%06d", n.Int64())
	if err := h.authSvc.SendSMS(req.Phone, code); err != nil {
		switch {
		case errors.Is(err, service.ErrSMSTooFrequent):
			response.Error(c, http.StatusTooManyRequests, errcode.ErrValidation, err.Error())
		case errors.Is(err, service.ErrPhoneOrCodeRequired):
			response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, "failed to send sms")
		}
		return
	}
	response.Success(c, gin.H{"status": "sent"})
}

// Profile 返回当前登录用户的基础身份信息。
func (h *UserHandler) Profile(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextKeyUserID) // M-01: C端使用 ContextKeyUserID
	if !exists {
		response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "not authenticated")
		return
	}

	response.Success(c, gin.H{"id": userID})
}
