package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/cruisebooking/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type cabinPricingVoyageStore interface {
	List(ctx context.Context) ([]domain.Voyage, error)
}

type cabinPricingCruiseStore interface {
	List(ctx context.Context, companyID int64, keyword string, status *int16, sortBy string, page, pageSize int) ([]domain.Cruise, int64, error)
}

// CabinPricingHandler 提供航次舱型价格管理端点。
type CabinPricingHandler struct {
	priceSvc  *service.VoyageCabinTypePriceService
	voyageSvc cabinPricingVoyageStore
	cruiseSvc cabinPricingCruiseStore
}

var shanghaiLocation = mustLoadLocation("Asia/Shanghai")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.FixedZone("CST", 8*60*60)
	}
	return loc
}

func NewCabinPricingHandler(priceSvc *service.VoyageCabinTypePriceService, voyageSvc cabinPricingVoyageStore, cruiseSvc cabinPricingCruiseStore) *CabinPricingHandler {
	return &CabinPricingHandler{priceSvc: priceSvc, voyageSvc: voyageSvc, cruiseSvc: cruiseSvc}
}

type BatchApplyPriceRequest struct {
	VoyageIDs            []int64 `json:"voyage_ids" binding:"required"`
	CabinTypeID          int64   `json:"cabin_type_id" binding:"required"`
	InventoryTotal       int     `json:"inventory_total" binding:"required"`
	SettlementPriceCents int64   `json:"settlement_price_cents" binding:"required"`
	SalePriceCents       int64   `json:"sale_price_cents" binding:"required"`
	EffectiveAt          string  `json:"effective_at"`
}

func (h *CabinPricingHandler) ListVoyages(c *gin.Context) {
	voyages, err := h.voyageSvc.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	companyID := queryInt64(c, "company_id", 0)
	cruiseID := queryInt64(c, "cruise_id", 0)
	departStart := c.Query("depart_start")
	departEnd := c.Query("depart_end")
	allowedCruiseIDs := map[int64]struct{}{}
	if companyID > 0 {
		cruises, _, err := h.cruiseSvc.List(c.Request.Context(), companyID, "", nil, "", 1, 10000)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
			return
		}
		for _, cruise := range cruises {
			allowedCruiseIDs[cruise.ID] = struct{}{}
		}
	}
	start, _ := parseDateOnly(departStart)
	end, _ := parseDateOnly(departEnd)
	filtered := make([]domain.Voyage, 0, len(voyages))
	for _, v := range voyages {
		if cruiseID > 0 && v.CruiseID != cruiseID {
			continue
		}
		if companyID > 0 {
			if _, ok := allowedCruiseIDs[v.CruiseID]; !ok {
				continue
			}
		}
		if !start.IsZero() && v.DepartDate.Before(start) {
			continue
		}
		if !end.IsZero() && v.DepartDate.After(end) {
			continue
		}
		filtered = append(filtered, v)
	}
	response.Success(c, gin.H{"list": filtered, "total": len(filtered)})
}

func (h *CabinPricingHandler) BatchApply(c *gin.Context) {
	var req BatchApplyPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if len(req.VoyageIDs) == 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "voyage_ids cannot be empty")
		return
	}
	effectiveAt := time.Now().In(shanghaiLocation)
	if req.EffectiveAt != "" {
		parsed, err := parseEffectiveAtInShanghai(req.EffectiveAt)
		if err != nil {
			response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid effective_at")
			return
		}
		effectiveAt = parsed
	}
	type applyFailure struct {
		VoyageID int64  `json:"voyage_id"`
		Reason   string `json:"reason"`
	}
	applied := 0
	failures := make([]applyFailure, 0)
	for _, voyageID := range req.VoyageIDs {
		v := &domain.VoyageCabinTypePriceVersion{
			VoyageID:             voyageID,
			CabinTypeID:          req.CabinTypeID,
			InventoryTotal:       req.InventoryTotal,
			SettlementPriceCents: req.SettlementPriceCents,
			SalePriceCents:       req.SalePriceCents,
			EffectiveAt:          effectiveAt,
		}
		if err := h.priceSvc.ApplyVersionAndRefreshCurrent(c.Request.Context(), v); err != nil {
			failures = append(failures, applyFailure{VoyageID: voyageID, Reason: err.Error()})
			continue
		}
		applied++
	}
	response.Success(c, gin.H{
		"applied": applied,
		"failed":  len(failures),
		"errors":  failures,
	})
}

func (h *CabinPricingHandler) History(c *gin.Context) {
	voyageID := queryInt64(c, "voyage_id", 0)
	cabinTypeID := queryInt64(c, "cabin_type_id", 0)
	if voyageID <= 0 || cabinTypeID <= 0 {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "voyage_id and cabin_type_id are required")
		return
	}
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 20)
	items, total, err := h.priceSvc.ListVersions(c.Request.Context(), voyageID, cabinTypeID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"list": items, "total": total})
}

func parseDateOnly(raw string) (time.Time, bool) {
	if raw == "" {
		return time.Time{}, false
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func parseEffectiveAtInShanghai(raw string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.In(shanghaiLocation), nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", raw, shanghaiLocation); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", raw, shanghaiLocation); err == nil {
		return t, nil
	}
	return time.Time{}, errors.New("invalid time format")
}
