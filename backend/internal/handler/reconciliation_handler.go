package handler

import (
	"net/http"
	"time"

	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type ReconciliationService interface {
	GenerateDailyReport(date time.Time) (interface{}, error)
}

type ReconciliationHandler struct {
	svc ReconciliationService
}

func NewReconciliationHandler(svc ReconciliationService) *ReconciliationHandler {
	return &ReconciliationHandler{svc: svc}
}

func (h *ReconciliationHandler) GenerateDailyReport(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		response.Error(c, http.StatusBadRequest, 400, "date is required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "invalid date format, use YYYY-MM-DD")
		return
	}

	report, err := h.svc.GenerateDailyReport(date)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, report)
}
