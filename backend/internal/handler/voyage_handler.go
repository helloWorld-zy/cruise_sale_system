package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// VoyageService 定义航次处理器所需的服务接口。
type VoyageService interface {
	List(ctx context.Context) ([]domain.Voyage, error)             // 查询航次列表
	Create(ctx context.Context, v *domain.Voyage) error            // 创建航次
	Update(ctx context.Context, v *domain.Voyage) error            // 更新航次
	GetByID(ctx context.Context, id int64) (*domain.Voyage, error) // 根据 ID 查询航次
	Delete(ctx context.Context, id int64) error                    // 删除航次
}

// VoyageHandler 处理 /admin/voyages 相关的 HTTP 端点。
// CRITICAL-03b + MEDIUM-04：实现完整的依赖注入和 CRUD 操作。
type VoyageHandler struct{ svc VoyageService }

// NewVoyageHandler 创建航次处理器实例。
func NewVoyageHandler(svc VoyageService) *VoyageHandler { return &VoyageHandler{svc: svc} }

type voyageItineraryPayload struct {
	DayNo             int    `json:"day_no" binding:"required,min=1"`
	StopIndex         int    `json:"stop_index" binding:"required,min=1"`
	City              string `json:"city" binding:"required,max=120"`
	Summary           string `json:"summary"`
	ETATime           string `json:"eta_time"`
	ETDTime           string `json:"etd_time"`
	HasBreakfast      bool   `json:"has_breakfast"`
	HasLunch          bool   `json:"has_lunch"`
	HasDinner         bool   `json:"has_dinner"`
	HasAccommodation  bool   `json:"has_accommodation"`
	AccommodationText string `json:"accommodation_text"`
}

type voyageUpsertPayload struct {
	CruiseID    int64                    `json:"cruise_id" binding:"required,min=1"`
	Code        string                   `json:"code" binding:"required,max=50"`
	ImageURL    string                   `json:"image_url" binding:"max=500"`
	BriefInfo   string                   `json:"brief_info" binding:"required,max=300"`
	DepartDate  time.Time                `json:"depart_date" binding:"required"`
	ReturnDate  time.Time                `json:"return_date" binding:"required"`
	Status      int16                    `json:"status"`
	Itineraries []voyageItineraryPayload `json:"itineraries" binding:"required,min=1,dive"`
}

// List 查询航次列表。
func (h *VoyageHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, list)
}

// Get 查询单条航次详情。
func (h *VoyageHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "voyage not found"})
		return
	}
	response.Success(c, item)
}

// Create 创建新的航次。
func (h *VoyageHandler) Create(c *gin.Context) {
	var req voyageUpsertPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v, err := buildVoyageFromPayload(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Create(c.Request.Context(), v); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, v)
}

// Update 更新指定的航次信息。
func (h *VoyageHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req voyageUpsertPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v, err := buildVoyageFromPayload(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v.ID = id
	if err := h.svc.Update(c.Request.Context(), v); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, v)
}

// Delete 删除指定的航次。
func (h *VoyageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondDeleteError(c, err, "voyage")
		return
	}
	c.Status(http.StatusNoContent)
}

func buildVoyageFromPayload(req voyageUpsertPayload) (*domain.Voyage, error) {
	if req.DepartDate.After(req.ReturnDate) {
		return nil, fmt.Errorf("depart_date must be before or equal to return_date")
	}
	v := &domain.Voyage{
		CruiseID:   req.CruiseID,
		Code:       strings.TrimSpace(req.Code),
		ImageURL:   strings.TrimSpace(req.ImageURL),
		BriefInfo:  strings.TrimSpace(req.BriefInfo),
		DepartDate: req.DepartDate,
		ReturnDate: req.ReturnDate,
		Status:     req.Status,
	}
	itineraries := make([]domain.VoyageItinerary, 0, len(req.Itineraries))
	seen := map[string]struct{}{}
	maxDay := 0
	for _, item := range req.Itineraries {
		key := fmt.Sprintf("%d-%d", item.DayNo, item.StopIndex)
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("duplicate itinerary day_no and stop_index: %s", key)
		}
		seen[key] = struct{}{}
		if item.DayNo > maxDay {
			maxDay = item.DayNo
		}
		itineraries = append(itineraries, domain.VoyageItinerary{
			DayNo:             item.DayNo,
			StopIndex:         item.StopIndex,
			City:              strings.TrimSpace(item.City),
			Summary:           strings.TrimSpace(item.Summary),
			ETATime:           normalizeTimeString(item.ETATime),
			ETDTime:           normalizeTimeString(item.ETDTime),
			HasBreakfast:      item.HasBreakfast,
			HasLunch:          item.HasLunch,
			HasDinner:         item.HasDinner,
			HasAccommodation:  item.HasAccommodation,
			AccommodationText: strings.TrimSpace(item.AccommodationText),
		})
	}
	for day := 1; day <= maxDay; day++ {
		if !hasItineraryDay(itineraries, day) {
			return nil, fmt.Errorf("itinerary day_no must be continuous from 1")
		}
	}
	v.Itineraries = itineraries
	return v, nil
}

func hasItineraryDay(itineraries []domain.VoyageItinerary, day int) bool {
	for _, item := range itineraries {
		if item.DayNo == day {
			return true
		}
	}
	return false
}

func normalizeTimeString(s string) *string {
	v := strings.TrimSpace(s)
	if v == "" {
		return nil
	}
	return &v
}
