package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type NotificationTemplateService interface {
	List(ctx context.Context) ([]domain.NotificationTemplate, error)
	Create(ctx context.Context, tpl *domain.NotificationTemplate) error
	Update(ctx context.Context, tpl *domain.NotificationTemplate) error
	Delete(ctx context.Context, id int64) error
}

type NotificationTemplateHandler struct {
	svc NotificationTemplateService
}

func NewNotificationTemplateHandler(svc NotificationTemplateService) *NotificationTemplateHandler {
	return &NotificationTemplateHandler{svc: svc}
}

func (h *NotificationTemplateHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *NotificationTemplateHandler) Create(c *gin.Context) {
	var req domain.NotificationTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	response.Success(c, req)
}

func (h *NotificationTemplateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	var req domain.NotificationTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	req.ID = id
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	response.Success(c, req)
}

func (h *NotificationTemplateHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
