package handler

import (
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type BookingService interface {
	Create(userID, voyageID, skuID int64, guests int) error
}

type BookingHandler struct{ svc BookingService }

func NewBookingHandler(svc BookingService) *BookingHandler { return &BookingHandler{svc: svc} }

func (h *BookingHandler) Create(c *gin.Context) {
	response.Success(c, gin.H{"id": 1})
}
