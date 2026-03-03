package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// CabinService 定义舱房处理器所需的服务接口。
// 包含 SKU 的 CRUD、库存管理和价格管理功能。
type CabinService interface {
	ListByVoyage(ctx context.Context, voyageID int64) ([]domain.CabinSKU, error)                 // 查询航次下的舱房列表
	FilteredList(ctx context.Context, f domain.CabinSKUFilter) ([]domain.CabinSKU, int64, error) // 高级筛选分页查询
	BatchUpdateStatus(ctx context.Context, ids []int64, status int16) error                      // 批量上下架
	GetByID(ctx context.Context, id int64) (*domain.CabinSKU, error)                             // 查询单个舱房 SKU
	Create(ctx context.Context, s *domain.CabinSKU) error                                        // 创建舱房 SKU
	Update(ctx context.Context, s *domain.CabinSKU) error                                        // 更新舱房 SKU
	Delete(ctx context.Context, id int64) error                                                  // 删除舱房 SKU
	GetInventory(ctx context.Context, skuID int64) (domain.CabinInventory, error)                // 查询库存
	GetAlerts(ctx context.Context) ([]domain.InventoryAlert, error)                              // 查询库存预警
	SetAlertThreshold(ctx context.Context, skuID int64, threshold int) error                     // 设置库存预警阈值
	AdjustInventory(ctx context.Context, skuID int64, delta int, reason string) error            // 调整库存
	ListPrices(ctx context.Context, skuID int64) ([]domain.CabinPrice, error)                    // 查询价格列表
	UpsertPrice(ctx context.Context, p *domain.CabinPrice) error                                 // 新增或更新价格
	BatchSetPrice(ctx context.Context, skuID int64, start, end time.Time, occupancy int, priceCents, childPriceCents, singleSupplementCents int64, priceType string) error
	GetCategoryTree(ctx context.Context) (interface{}, error) // 获取邮轮→航线→舱型分类树
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

// CabinBatchStatusRequest 表示批量更新舱位状态请求。
type CabinBatchStatusRequest struct {
	IDs    []int64 `json:"ids" binding:"required"` // 舱位 ID 列表
	Status int16   `json:"status"`                 // 目标状态
}

// CabinAlertThresholdRequest 表示库存预警阈值请求。
type CabinAlertThresholdRequest struct {
	Threshold int `json:"threshold"` // 预警阈值
}

// CabinAdjustInventoryRequest 表示库存调整请求。
type CabinAdjustInventoryRequest struct {
	Delta  int    `json:"delta" binding:"required"`  // 调整量（正数增加，负数减少）
	Reason string `json:"reason" binding:"required"` // 调整原因
}

// CabinBatchPriceRequest 表示按日期区间批量设置价格请求。
type CabinBatchPriceRequest struct {
	StartDate             string `json:"start_date" binding:"required"`
	EndDate               string `json:"end_date" binding:"required"`
	Occupancy             int    `json:"occupancy" binding:"required,gt=0"`
	PriceCents            int64  `json:"price_cents" binding:"required,gte=0"`
	ChildPriceCents       int64  `json:"child_price_cents"`
	SingleSupplementCents int64  `json:"single_supplement_cents"`
	PriceType             string `json:"price_type"`
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

// FilteredList godoc
// @Summary Filter cabin sku list
// @Tags Cabin
// @Security BearerAuth
// @Produce json
// @Param voyage_id query int false "Voyage ID"
// @Param cabin_type_id query int false "Cabin type ID"
// @Param status query int false "Status, 0 means all"
// @Param keyword query string false "Keyword"
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins [get]
// FilteredList 按条件分页查询舱位列表。
func (h *CabinHandler) FilteredList(c *gin.Context) {
	var statusPtr *int16
	if s := c.Query("status"); s != "" {
		v := int16(queryInt(c, "status", 0))
		statusPtr = &v
	}
	filter := domain.CabinSKUFilter{
		VoyageID:    queryInt64(c, "voyage_id", 0),
		CabinTypeID: queryInt64(c, "cabin_type_id", 0),
		Status:      statusPtr,
		Keyword:     c.Query("keyword"),
		Page:        queryInt(c, "page", 1),
		PageSize:    queryInt(c, "page_size", 10),
	}
	list, total, err := h.svc.FilteredList(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

// BatchUpdateStatus godoc
// @Summary Batch update cabin sku status
// @Tags Cabin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CabinBatchStatusRequest true "Batch status payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/batch-status [put]
// BatchUpdateStatus 批量更新舱位状态。
func (h *CabinHandler) BatchUpdateStatus(c *gin.Context) {
	var req CabinBatchStatusRequest
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
	log.Printf("audit bulk update cabins count=%d status=%d", len(req.IDs), req.Status)
	if err := h.svc.BatchUpdateStatus(c.Request.Context(), req.IDs, req.Status); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, gin.H{"updated": len(req.IDs)})
}

// GetAlerts godoc
// @Summary List inventory alerts
// @Tags Cabin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/alerts [get]
// GetAlerts 查询库存预警列表。
func (h *CabinHandler) GetAlerts(c *gin.Context) {
	alerts, err := h.svc.GetAlerts(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, alerts)
}

// SetAlertThreshold godoc
// @Summary Set cabin inventory alert threshold
// @Tags Cabin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Cabin SKU ID"
// @Param body body CabinAlertThresholdRequest true "Alert threshold payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/{id}/alert-threshold [put]
// SetAlertThreshold 设置舱位库存预警阈值。
func (h *CabinHandler) SetAlertThreshold(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	var req CabinAlertThresholdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if req.Threshold < 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "threshold must be >= 0")
		return
	}
	if err := h.svc.SetAlertThreshold(c.Request.Context(), id, req.Threshold); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, gin.H{"id": id, "threshold": req.Threshold})
}

// Get godoc
// @Summary Get cabin sku detail
// @Tags Cabin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Cabin SKU ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/admin/cabins/{id} [get]
// Get 查询单个舱房 SKU 详情。
func (h *CabinHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "cabin not found")
		return
	}
	response.Success(c, item)
}

// Create godoc
// @Summary Create cabin sku
// @Tags Cabin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body domain.CabinSKU true "Cabin SKU payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins [post]
// Create 创建新的舱房 SKU。
func (h *CabinHandler) Create(c *gin.Context) {
	var req domain.CabinSKU
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}
	h.indexCabinDocument(req)
	response.Success(c, req)
}

// Update godoc
// @Summary Update cabin sku
// @Tags Cabin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Cabin SKU ID"
// @Param body body domain.CabinSKU true "Cabin SKU payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/{id} [put]
// Update 更新指定的舱房 SKU。
func (h *CabinHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	var req domain.CabinSKU
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
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

// Delete godoc
// @Summary Delete cabin sku
// @Tags Cabin
// @Security BearerAuth
// @Param id path int true "Cabin SKU ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/{id} [delete]
// Delete 删除指定的舱房 SKU。
func (h *CabinHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondDeleteError(c, err, "cabin")
		return
	}
	c.Status(http.StatusNoContent)
}

// GetInventory godoc
// @Summary Get cabin inventory
// @Tags Cabin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Cabin SKU ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/{id}/inventory [get]
// GetInventory 查询指定舱房 SKU 的库存信息。
// GET /admin/cabins/:id/inventory
func (h *CabinHandler) GetInventory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	inv, err := h.svc.GetInventory(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, inv)
}

// AdjustInventory godoc
// @Summary Adjust cabin inventory
// @Tags Cabin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Cabin SKU ID"
// @Param body body CabinAdjustInventoryRequest true "Inventory adjustment payload"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/{id}/inventory/adjust [post]
// AdjustInventory 调整指定舱房 SKU 的库存数量。
// POST /admin/cabins/:id/inventory/adjust
func (h *CabinHandler) AdjustInventory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	// 请求体：库存调整量和调整原因
	var req CabinAdjustInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if err := h.svc.AdjustInventory(c.Request.Context(), id, req.Delta, req.Reason); err != nil {
		response.InternalError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListPrices godoc
// @Summary List cabin price calendar
// @Tags Cabin
// @Security BearerAuth
// @Produce json
// @Param id path int true "Cabin SKU ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/{id}/prices [get]
// ListPrices 查询指定舱房 SKU 的价格日历。
// GET /admin/cabins/:id/prices
func (h *CabinHandler) ListPrices(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	prices, err := h.svc.ListPrices(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, prices)
}

// UpsertPrice godoc
// @Summary Create or update cabin price
// @Tags Cabin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Cabin SKU ID"
// @Param body body domain.CabinPrice true "Cabin price payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/{id}/prices [post]
// UpsertPrice 新增或更新指定舱房 SKU 的价格记录。
// POST /admin/cabins/:id/prices
func (h *CabinHandler) UpsertPrice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	var req domain.CabinPrice
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	req.CabinSKUID = id
	if err := h.svc.UpsertPrice(c.Request.Context(), &req); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, req)
}

// BatchSetPrice 按日期区间批量设置指定舱位 SKU 的价格。
func (h *CabinHandler) BatchSetPrice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	var req CabinBatchPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid start_date")
		return
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid end_date")
		return
	}
	if end.Before(start) {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "end_date must be >= start_date")
		return
	}
	if req.PriceType == "" {
		req.PriceType = "base"
	}
	if err := h.svc.BatchSetPrice(c.Request.Context(), id, start, end, req.Occupancy, req.PriceCents, req.ChildPriceCents, req.SingleSupplementCents, req.PriceType); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, gin.H{"id": id, "start_date": req.StartDate, "end_date": req.EndDate})
}

// CategoryTree godoc
// @Summary Get cabin category tree (cruise -> route -> cabin type)
// @Tags Cabin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/admin/cabins/category-tree [get]
// CategoryTree 获取邮轮→航线→舱型三级分类树。
func (h *CabinHandler) CategoryTree(c *gin.Context) {
	tree, err := h.svc.GetCategoryTree(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, tree)
}
