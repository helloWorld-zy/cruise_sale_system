package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// RouteService 定义航线处理器所需的服务接口。
// MEDIUM-05 修复：使用 domain.Route 代替 []interface{}。
type RouteService interface {
	List(ctx context.Context) ([]domain.Route, error)             // 查询所有航线
	Create(ctx context.Context, r *domain.Route) error            // 创建航线
	Update(ctx context.Context, r *domain.Route) error            // 更新航线
	GetByID(ctx context.Context, id int64) (*domain.Route, error) // 根据 ID 查询航线
	Delete(ctx context.Context, id int64) error                   // 删除航线
}

// RouteHandler 处理 /admin/routes 相关的 HTTP 端点。
// CRITICAL-03b + HIGH-03：所有错误均已正确处理，不再被静默忽略。
type RouteHandler struct{ svc RouteService }

// NewRouteHandler 创建航线处理器实例。
func NewRouteHandler(svc RouteService) *RouteHandler { return &RouteHandler{svc: svc} }

// List 查询所有航线。
func (h *RouteHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, list)
}

// Get 查询单条航线详情。
func (h *RouteHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
		return
	}
	response.Success(c, item)
}

// Create 创建新的航线。
func (h *RouteHandler) Create(c *gin.Context) {
	var req domain.Route
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

// Update 更新指定的航线信息。
func (h *RouteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.Route
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

// Delete 删除指定的航线。
func (h *RouteHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondDeleteError(c, err, "route")
		return
	}
	c.Status(http.StatusNoContent)
}
