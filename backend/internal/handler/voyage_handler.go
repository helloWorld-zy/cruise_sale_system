package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// VoyageService 定义航次处理器所需的服务接口。
type VoyageService interface {
	ListByRoute(ctx context.Context, routeID int64) ([]domain.Voyage, error) // 查询某航线下的航次列表
	Create(ctx context.Context, v *domain.Voyage) error                      // 创建航次
	Update(ctx context.Context, v *domain.Voyage) error                      // 更新航次
	Delete(ctx context.Context, id int64) error                              // 删除航次
}

// VoyageHandler 处理 /admin/voyages 相关的 HTTP 端点。
// CRITICAL-03b + MEDIUM-04：实现完整的依赖注入和 CRUD 操作。
type VoyageHandler struct{ svc VoyageService }

// NewVoyageHandler 创建航次处理器实例。
func NewVoyageHandler(svc VoyageService) *VoyageHandler { return &VoyageHandler{svc: svc} }

// List 查询指定航线下的航次列表。
func (h *VoyageHandler) List(c *gin.Context) {
	routeID, _ := strconv.ParseInt(c.Query("route_id"), 10, 64)
	list, err := h.svc.ListByRoute(c.Request.Context(), routeID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, list)
}

// Create 创建新的航次。
func (h *VoyageHandler) Create(c *gin.Context) {
	var req domain.Voyage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, req)
}

// Update 更新指定的航次信息。
func (h *VoyageHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.Voyage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, req)
}

// Delete 删除指定的航次。
func (h *VoyageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.InternalError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
