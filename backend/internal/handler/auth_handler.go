package handler

import (
	"net/http"
	"time"

	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles login and token-related endpoints.
type AuthHandler struct {
	authSvc *service.AuthService
}

// NewAuthHandler creates an AuthHandler with injected AuthService.
func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// LoginRequest is the payload for POST /api/v1/admin/auth/login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse contains the issued JWT token and profile info.
type LoginResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
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
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

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
func (h *AuthHandler) GetProfile(c *gin.Context) {
	staffIDStr, exists := c.Get("staffID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, errcode.ErrUnauthorized, "not authenticated")
		return
	}

	staff, err := h.authSvc.GetProfile(c.Request.Context(), staffIDStr.(string))
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "staff not found")
		return
	}

	response.Success(c, staff)
}
