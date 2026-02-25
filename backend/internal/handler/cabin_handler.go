package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// CabinService 定义舱房处理器所需的服务接口。
// 包含 SKU 的 CRUD、库存管理和价格管理功能。
type CabinService interface {
	ListByVoyage(ctx context.Context, voyageID int64) ([]domain.CabinSKU, error)      // 查询航次下的舱房列表
	Create(ctx context.Context, s *domain.CabinSKU) error                             // 创建舱房 SKU
	Update(ctx context.Context, s *domain.CabinSKU) error                             // 更新舱房 SKU
	Delete(ctx context.Context, id int64) error                                       // 删除舱房 SKU
	GetInventory(ctx context.Context, skuID int64) (domain.CabinInventory, error)     // 查询库存
	AdjustInventory(ctx context.Context, skuID int64, delta int, reason string) error // 调整库存
	ListPrices(ctx context.Context, skuID int64) ([]domain.CabinPrice, error)         // 查询价格列表
	UpsertPrice(ctx context.Context, p *domain.CabinPrice) error                      // 新增或更新价格
}

// CabinIndexer 定义舱位文档索引能力。
type CabinIndexer interface {
	IndexCabin(doc interface{}) error
}

// CabinIndexRetryQueue 定义索引失败重试队列能力。
type CabinIndexRetryQueue interface {
	Enqueue(doc interface{})
}

// CabinHandler 处理 /admin/cabins 相关的 HTTP 端点。
// 提供完整的 CRUD 以及库存和定价子资源的管理功能。
type CabinHandler struct {
	svc        CabinService
	indexer    CabinIndexer
	retryQueue CabinIndexRetryQueue
}

// NewCabinHandler 创建舱房处理器实例。
func NewCabinHandler(svc CabinService) *CabinHandler { return &CabinHandler{svc: svc} }

// NewCabinHandlerWithIndexing 创建带索引与重试能力的舱房处理器实例。
func NewCabinHandlerWithIndexing(svc CabinService, indexer CabinIndexer, retryQueue CabinIndexRetryQueue) *CabinHandler {
	return &CabinHandler{svc: svc, indexer: indexer, retryQueue: retryQueue}
}

// indexCabinDocument 将舱位变更同步到搜索索引，失败时进入重试队列。
func (h *CabinHandler) indexCabinDocument(doc domain.CabinSKU) {
	if h.indexer == nil {
		return
	}
	if err := h.indexer.IndexCabin(doc); err != nil && h.retryQueue != nil {
		h.retryQueue.Enqueue(doc)
	}
}

// List 查询指定航次下的所有舱房 SKU。
func (h *CabinHandler) List(c *gin.Context) {
	voyageID, _ := strconv.ParseInt(c.Query("voyage_id"), 10, 64)
	list, err := h.svc.ListByVoyage(c.Request.Context(), voyageID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, list)
}

// Create 创建新的舱房 SKU。
func (h *CabinHandler) Create(c *gin.Context) {
	var req domain.CabinSKU
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}
	h.indexCabinDocument(req)
	response.Success(c, req)
}

// Update 更新指定的舱房 SKU。
func (h *CabinHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.CabinSKU
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}
	h.indexCabinDocument(req)
	response.Success(c, req)
}

// Delete 删除指定的舱房 SKU。
func (h *CabinHandler) Delete(c *gin.Context) {
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

// GetInventory 查询指定舱房 SKU 的库存信息。
// GET /admin/cabins/:id/inventory
func (h *CabinHandler) GetInventory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	inv, err := h.svc.GetInventory(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, inv)
}

// AdjustInventory 调整指定舱房 SKU 的库存数量。
// POST /admin/cabins/:id/inventory/adjust
func (h *CabinHandler) AdjustInventory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	// 请求体：库存调整量和调整原因
	var req struct {
		Delta  int    `json:"delta" binding:"required"`  // 调整量（正数为增加，负数为减少）
		Reason string `json:"reason" binding:"required"` // 调整原因说明
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AdjustInventory(c.Request.Context(), id, req.Delta, req.Reason); err != nil {
		response.InternalError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListPrices 查询指定舱房 SKU 的价格日历。
// GET /admin/cabins/:id/prices
func (h *CabinHandler) ListPrices(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	prices, err := h.svc.ListPrices(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, prices)
}

// UpsertPrice 新增或更新指定舱房 SKU 的价格记录。
// POST /admin/cabins/:id/prices
func (h *CabinHandler) UpsertPrice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.CabinPrice
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.CabinSKUID = id
	if err := h.svc.UpsertPrice(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, req)
}
