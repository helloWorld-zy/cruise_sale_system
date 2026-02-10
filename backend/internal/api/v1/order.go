package v1

import (
	"cruise_booking_system/internal/core"
	"cruise_booking_system/internal/model"
	"cruise_booking_system/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderSvc   *core.OrderService
	paymentSvc *core.PaymentService
}

func NewOrderHandler(orderSvc *core.OrderService, paymentSvc *core.PaymentService) *OrderHandler {
	return &OrderHandler{
		orderSvc:   orderSvc,
		paymentSvc: paymentSvc,
	}
}

type CreateOrderRequest struct {
	VoyageID    string `json:"voyage_id" binding:"required"`
	CabinTypeID string `json:"cabin_type_id" binding:"required"`
	CabinID     string `json:"cabin_id"`
	Passengers  []struct {
		NameCn    string `json:"name_cn"`
		NameEn    string `json:"name_en"`
		DocType   string `json:"doc_type"`
		DocNumber string `json:"doc_number"`
	} `json:"passengers"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	userIDStr := c.GetString("userID")
	// userID, _ := uuid.Parse(userIDStr) // Handle error in real app
	// Mocking user ID if not present (dev mode) or assuming middleware set it
	if userIDStr == "" {
		userIDStr = uuid.New().String()
	}
	userID, _ := uuid.Parse(userIDStr)

	voyageID, _ := uuid.Parse(req.VoyageID)
	cabinTypeID, _ := uuid.Parse(req.CabinTypeID)
	var cabinID *uuid.UUID
	if req.CabinID != "" {
		cid, _ := uuid.Parse(req.CabinID)
		cabinID = &cid
	}

	order := &model.Order{
		VoyageID: voyageID,
		Items: []model.OrderItem{
			{
				CabinTypeID: cabinTypeID,
				CabinID:     cabinID,
				// Price would be fetched from PriceRule
			},
		},
	}
	// Map passengers...

	createdOrder, err := h.orderSvc.CreateOrder(c.Request.Context(), userID, order)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	payURL := h.paymentSvc.GetPaymentURL(createdOrder.OrderNo, createdOrder.TotalAmount)

	response.Success(c, gin.H{
		"order_id":    createdOrder.ID,
		"payment_url": payURL,
	})
}

func (h *OrderHandler) ListMine(c *gin.Context) {
	userIDStr := c.GetString("userID")
	if userIDStr == "" {
		// For dev/test without auth
		userIDStr = uuid.New().String() 
	}
	userID, _ := uuid.Parse(userIDStr)

	orders, err := h.orderSvc.ListUserOrders(c.Request.Context(), userID)
	if err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, orders)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	userIDStr := c.GetString("userID")
	userID, _ := uuid.Parse(userIDStr)

	if err := h.orderSvc.CancelOrder(c.Request.Context(), userID, id); err != nil {
		response.ServerError(c, err)
		return
	}
	response.Success(c, nil)
}
