package handler

import (
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// CabinTypeCategoryHandler 处理舱型大类字典管理端点。
type CabinTypeCategoryHandler struct {
	svc *service.CabinTypeCategoryService
}

func NewCabinTypeCategoryHandler(svc *service.CabinTypeCategoryService) *CabinTypeCategoryHandler {
	return &CabinTypeCategoryHandler{svc: svc}
}

type CabinTypeCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Status    int16  `json:"status"`
	SortOrder int    `json:"sort_order"`
}

func (h *CabinTypeCategoryHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *CabinTypeCategoryHandler) Create(c *gin.Context) {
	var req CabinTypeCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	item := &domain.CabinTypeCategory{
		Name:      req.Name,
		Code:      req.Code,
		Status:    req.Status,
		SortOrder: req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), item); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *CabinTypeCategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "cabin type category not found")
		return
	}
	var req CabinTypeCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	item.Name = req.Name
	item.Code = req.Code
	item.Status = req.Status
	item.SortOrder = req.SortOrder
	if err := h.svc.Update(c.Request.Context(), item); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *CabinTypeCategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, nil)
}
