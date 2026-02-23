package handler

import (
	"net/http"
	"time"

	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler 处理登录和令牌相关的 HTTP 端点。
type AuthHandler struct {
	authSvc *service.AuthService // 认证服务
}

// NewAuthHandler 创建认证处理器，通过依赖注入传入认证服务。
func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// LoginRequest 是 POST /api/v1/admin/auth/login 的请求体结构。
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 登录用户名
	Password string `json:"password" binding:"required"` // 登录密码
}

// LoginResponse 包含签发的 JWT 令牌和过期时间。
type LoginResponse struct {
	Token    string    `json:"token"`     // JWT 令牌字符串
	ExpireAt time.Time `json:"expire_at"` // 令牌过期时间
}

// Login godoc
// @Summary Admin login
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/admin/auth/login [post]
// Login 处理管理员登录请求。
// 验证请求参数 → 调用认证服务校验凭据 → 返回 JWT 令牌。
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	// 调用认证服务进行登录验证
	token, expireAt, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "invalid credentials")
		return
	}

	response.Success(c, LoginResponse{Token: token, ExpireAt: expireAt})
}

// GetProfile godoc
// @Summary Get current staff profile
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/admin/auth/profile [get]
// GetProfile 获取当前登录员工的个人信息。
// 从 gin 上下文中读取 JWT 中间件设置的 staffID。
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 从上下文获取已认证的员工 ID
	staffIDStr, exists := c.Get("staffID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "not authenticated")
		return
	}

	// 查询员工信息
	staff, err := h.authSvc.GetProfile(c.Request.Context(), staffIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "staff not found")
		return
	}

	response.Success(c, staff)
}
