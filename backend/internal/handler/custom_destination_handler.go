package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// CustomDestinationHandler 处理自定义目的地的 CRUD 端点。
type CustomDestinationHandler struct {
	svc *service.CustomDestinationService
}

// NewCustomDestinationHandler 创建自定义目的地处理器。
func NewCustomDestinationHandler(svc *service.CustomDestinationService) *CustomDestinationHandler {
	return &CustomDestinationHandler{svc: svc}
}

// CustomDestinationRequest 创建/更新自定义目的地的请求体。
type CustomDestinationRequest struct {
	Name        string   `json:"name" binding:"required"`    // 目的地名称（必填）
	Country     string   `json:"country" binding:"required"` // 所属国家（必填）
	Latitude    *float64 `json:"latitude"`                   // 纬度
	Longitude   *float64 `json:"longitude"`                  // 经度
	Keywords    string   `json:"keywords"`                   // 搜索关键词（逗号分隔）
	Description string   `json:"description"`                // 备注描述
	Status      int16    `json:"status"`                     // 启用状态
	SortOrder   int      `json:"sort_order"`                 // 排序权重
}

// List 查询所有自定义目的地。
func (h *CustomDestinationHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, items)
}

// Get 查询指定自定义目的地详情。
func (h *CustomDestinationHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "destination not found")
		return
	}
	response.Success(c, item)
}

// Create 创建新的自定义目的地。
func (h *CustomDestinationHandler) Create(c *gin.Context) {
	var req CustomDestinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if req.Latitude == nil || req.Longitude == nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "latitude and longitude are required")
		return
	}
	dest := &domain.CustomDestination{
		Name:        req.Name,
		Country:     req.Country,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Keywords:    req.Keywords,
		Description: req.Description,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), dest); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, dest)
}

// Update 更新指定自定义目的地。
func (h *CustomDestinationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	existing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "destination not found")
		return
	}
	var req CustomDestinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if req.Latitude == nil || req.Longitude == nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "latitude and longitude are required")
		return
	}
	existing.Name = req.Name
	existing.Country = req.Country
	existing.Latitude = req.Latitude
	existing.Longitude = req.Longitude
	existing.Keywords = req.Keywords
	existing.Description = req.Description
	existing.Status = req.Status
	existing.SortOrder = req.SortOrder
	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, existing)
}

// Delete 删除指定的自定义目的地。
func (h *CustomDestinationHandler) Delete(c *gin.Context) {
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

// ExportCSV 导出港口城市词典 CSV。
func (h *CustomDestinationHandler) ExportCSV(c *gin.Context) {
	data, err := h.svc.ExportCSV(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	filename := fmt.Sprintf("port_city_dictionary_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Data(http.StatusOK, "text/csv", data)
}

// ImportCSV 导入港口城市词典 CSV。
func (h *CustomDestinationHandler) ImportCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "csv file is required")
		return
	}
	opened, err := file.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "failed to read csv file")
		return
	}
	defer opened.Close()
	summary, err := h.svc.ImportCSV(c.Request.Context(), opened)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	response.Success(c, summary)
}
