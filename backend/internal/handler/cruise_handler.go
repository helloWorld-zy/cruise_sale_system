package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// CruiseHandler 处理邮轮的 CRUD 端点。
// CR-05：通过 CruiseService 实现依赖注入。
type CruiseHandler struct {
	svc *service.CruiseService // 邮轮服务
}

// NewCruiseHandler 创建邮轮处理器，通过依赖注入传入服务。
func NewCruiseHandler(svc *service.CruiseService) *CruiseHandler {
	return &CruiseHandler{svc: svc}
}

// CruiseRequest 是创建/更新邮轮的请求体结构。
type CruiseRequest struct {
	CompanyID         int64   `json:"company_id" binding:"required"` // 所属公司 ID（必填）
	Name              string  `json:"name" binding:"required"`       // 邮轮名称（必填）
	EnglishName       string  `json:"english_name"`                  // 英文名称
	Code              string  `json:"code"`                          // 邮轮代码/英文简称
	CrewCount         int     `json:"crew_count"`                    // 船员人数
	BuildYear         int     `json:"build_year"`                    // 建造年份
	RefurbishYear     int     `json:"refurbish_year"`                // 最近翻新年份
	Tonnage           float64 `json:"tonnage"`                       // 吨位
	PassengerCapacity int     `json:"passenger_capacity"`            // 最大载客量
	RoomCount         int     `json:"room_count"`                    // 舱房总数
	Length            float64 `json:"length"`                        // 船体长度（米）
	Width             float64 `json:"width"`                         // 船体宽度（米）
	DeckCount         int     `json:"deck_count"`                    // 甲板层数
	Description       string  `json:"description"`                   // 描述
	SortOrder         int     `json:"sort_order"`                    // 排序权重
}

// CruiseBatchStatusRequest 是邮轮批量更新状态的请求结构。
type CruiseBatchStatusRequest struct {
	IDs    []int64 `json:"ids" binding:"required"` // 需要更新状态的邮轮 ID 列表
	Status int16   `json:"status"`                 // 目标状态
}

// List godoc
// @Summary List cruises
// @Tags Cruise
// @Security BearerAuth
// @Produce json
// @Param company_id query int false "Filter by company"
// @Param keyword query string false "Keyword on cruise name"
// @Param status query int false "Filter by status, 0 means all"
// @Param sort_by query string false "Sort by tonnage/name"
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cruises [get]
// List 分页查询邮轮列表，可按公司 ID 过滤。
func (h *CruiseHandler) List(c *gin.Context) {
	companyID := queryInt64(c, "company_id", 0)
	keyword := c.Query("keyword")
	var status *int16
	if s := c.Query("status"); s != "" {
		v := int16(queryInt(c, "status", 0))
		status = &v
	}
	sortBy := c.Query("sort_by")
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 10)

	items, total, err := h.svc.List(c.Request.Context(), companyID, keyword, status, sortBy, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"list": items, "total": total})
}

// Get 查询单个邮轮详情。
func (h *CruiseHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}

	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "cruise not found")
		return
	}

	response.Success(c, item)
}

// Create godoc
// @Summary Create cruise
// @Tags Cruise
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CruiseRequest true "Cruise data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/admin/cruises [post]
// Create 创建新的邮轮，会先验证所属公司是否存在。
func (h *CruiseHandler) Create(c *gin.Context) {
	var req CruiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	cruise := &domain.Cruise{
		CompanyID:         req.CompanyID,
		Name:              req.Name,
		EnglishName:       req.EnglishName,
		Code:              req.Code,
		CrewCount:         req.CrewCount,
		BuildYear:         req.BuildYear,
		RefurbishYear:     req.RefurbishYear,
		Tonnage:           req.Tonnage,
		PassengerCapacity: req.PassengerCapacity,
		RoomCount:         req.RoomCount,
		Length:            req.Length,
		Width:             req.Width,
		DeckCount:         req.DeckCount,
		Description:       req.Description,
		SortOrder:         req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), cruise); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrBadRequest, err.Error())
		return
	}
	response.Success(c, cruise)
}

// Update godoc
// @Summary Update cruise
// @Tags Cruise
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Cruise ID"
// @Param body body CruiseRequest true "Cruise data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cruises/{id} [put]
// Update 更新指定的邮轮信息。
func (h *CruiseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	// 查询现有邮轮数据
	existing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "cruise not found")
		return
	}
	var req CruiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	// 更新字段
	existing.Name = req.Name
	existing.EnglishName = req.EnglishName
	existing.Code = req.Code
	existing.CrewCount = req.CrewCount
	existing.BuildYear = req.BuildYear
	existing.RefurbishYear = req.RefurbishYear
	existing.Tonnage = req.Tonnage
	existing.PassengerCapacity = req.PassengerCapacity
	existing.RoomCount = req.RoomCount
	existing.Length = req.Length
	existing.Width = req.Width
	existing.DeckCount = req.DeckCount
	existing.Description = req.Description
	existing.SortOrder = req.SortOrder

	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, existing)
}

// Delete godoc
// @Summary Delete cruise
// @Tags Cruise
// @Security BearerAuth
// @Param id path int true "Cruise ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cruises/{id} [delete]
// Delete 删除指定的邮轮。若邮轮下仍有舱房类型则返回冲突错误。
func (h *CruiseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrCruiseHasCabins) {
			response.Error(c, http.StatusConflict, errcode.ErrCruiseHasCabins, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, nil)
}

// BatchUpdateStatus godoc
// @Summary Batch update cruise status
// @Tags Cruise
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CruiseBatchStatusRequest true "Batch status payload"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/cruises/batch-status [put]
// BatchUpdateStatus 批量更新邮轮状态。
func (h *CruiseHandler) BatchUpdateStatus(c *gin.Context) {
	var req CruiseBatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if len(req.IDs) == 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "ids cannot be empty")
		return
	}
	if len(req.IDs) > maxBatchUpdateSize {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "batch size exceeds limit")
		return
	}
	log.Printf("audit bulk update cruises count=%d status=%d", len(req.IDs), req.Status)
	if err := h.svc.BatchUpdateStatus(c.Request.Context(), req.IDs, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"updated": len(req.IDs)})
}
