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

// CompanyHandler 处理邮轮公司的 CRUD 端点。
// CR-05：通过 CompanyService 实现依赖注入。
type CompanyHandler struct {
	svc *service.CompanyService // 邮轮公司服务
}

// NewCompanyHandler 创建公司处理器，通过依赖注入传入服务。
func NewCompanyHandler(svc *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{svc: svc}
}

// CompanyRequest 是创建/更新公司的请求体结构。
type CompanyRequest struct {
	Name        string `json:"name" binding:"required"` // 公司名称（必填）
	EnglishName string `json:"english_name"`            // 英文名称
	Description string `json:"description"`             // 公司描述
	LogoURL     string `json:"logo_url"`                // Logo 图片地址
	SortOrder   int    `json:"sort_order"`              // 排序权重
}

// List godoc
// @Summary List cruise companies
// @Tags Company
// @Security BearerAuth
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/admin/companies [get]
// List 分页查询邮轮公司列表，支持关键词搜索。
func (h *CompanyHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 10)

	items, total, err := h.svc.List(c.Request.Context(), keyword, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"list": items, "total": total})
}

// Create godoc
// @Summary Create a cruise company
// @Tags Company
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CompanyRequest true "Company data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/companies [post]
// Create 创建新的邮轮公司。
func (h *CompanyHandler) Create(c *gin.Context) {
	var req CompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	company := &domain.CruiseCompany{
		Name:        req.Name,
		EnglishName: req.EnglishName,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		SortOrder:   req.SortOrder,
	}
	if err := h.svc.Create(c.Request.Context(), company); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, company)
}

// Update godoc
// @Summary Update a cruise company
// @Tags Company
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Company ID"
// @Param body body CompanyRequest true "Company data"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/companies/{id} [put]
// Update 更新指定的邮轮公司信息。
func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	// 查询现有公司数据
	existing, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "company not found")
		return
	}
	var req CompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	// 更新字段
	existing.Name = req.Name
	existing.EnglishName = req.EnglishName
	existing.Description = req.Description
	existing.LogoURL = req.LogoURL
	existing.SortOrder = req.SortOrder
	if err := h.svc.Update(c.Request.Context(), existing); err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, existing)
}

// Delete godoc
// @Summary Delete a cruise company
// @Tags Company
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/companies/{id} [delete]
// Delete 删除指定的邮轮公司。若公司下仍有邮轮则返回冲突错误。
func (h *CompanyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if err.Error() == "company has cruises" {
			response.Error(c, http.StatusConflict, errcode.ErrCompanyHasCruises, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, nil)
}
