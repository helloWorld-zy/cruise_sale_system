package handler

import (
	"net/http"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type ShopInfoService interface {
	Get(ctx interface{}) (*domain.ShopInfo, error)
	Update(ctx interface{}, info *domain.ShopInfo) error
}

type ShopInfoHandler struct {
	svc ShopInfoService
}

func NewShopInfoHandler(svc ShopInfoService) *ShopInfoHandler {
	return &ShopInfoHandler{svc: svc}
}

func (h *ShopInfoHandler) Get(c *gin.Context) {
	info, err := h.svc.Get(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if info == nil {
		info = &domain.ShopInfo{}
	}
	response.Success(c, info)
}

func (h *ShopInfoHandler) Update(c *gin.Context) {
	var req domain.ShopInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, req)
}
