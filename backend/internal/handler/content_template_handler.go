package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type ContentTemplateService interface {
	List(ctx context.Context, kind domain.ContentTemplateKind) ([]domain.ContentTemplate, error)
	GetByID(ctx context.Context, id int64) (*domain.ContentTemplate, error)
	Create(ctx context.Context, tpl *domain.ContentTemplate) error
	Update(ctx context.Context, tpl *domain.ContentTemplate) error
	Delete(ctx context.Context, id int64) error
}

type ContentTemplateHandler struct {
	svc ContentTemplateService
}

type contentTemplatePayload struct {
	Name    string                     `json:"name" binding:"required,max=120"`
	Kind    domain.ContentTemplateKind `json:"kind" binding:"required"`
	Status  int16                      `json:"status"`
	Content json.RawMessage            `json:"content" binding:"required"`
}

func NewContentTemplateHandler(svc ContentTemplateService) *ContentTemplateHandler {
	return &ContentTemplateHandler{svc: svc}
}

func (h *ContentTemplateHandler) List(c *gin.Context) {
	kind := domain.ContentTemplateKind(c.Query("kind"))
	list, err := h.svc.List(c.Request.Context(), kind)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, toContentTemplateResponses(list))
}

func (h *ContentTemplateHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	tpl, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	if tpl == nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "template not found")
		return
	}
	response.Success(c, toContentTemplateResponse(*tpl))
}

func (h *ContentTemplateHandler) Create(c *gin.Context) {
	var req contentTemplatePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	tpl := &domain.ContentTemplate{
		Name:        req.Name,
		Kind:        req.Kind,
		Status:      req.Status,
		ContentJSON: normalizeJSONBytes(req.Content),
	}
	if err := h.svc.Create(c.Request.Context(), tpl); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	response.Success(c, toContentTemplateResponse(*tpl))
}

func (h *ContentTemplateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	var req contentTemplatePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	tpl := &domain.ContentTemplate{
		ID:          id,
		Name:        req.Name,
		Kind:        req.Kind,
		Status:      req.Status,
		ContentJSON: normalizeJSONBytes(req.Content),
	}
	if err := h.svc.Update(c.Request.Context(), tpl); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	response.Success(c, toContentTemplateResponse(*tpl))
}

func (h *ContentTemplateHandler) Delete(c *gin.Context) {
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

func toContentTemplateResponses(list []domain.ContentTemplate) []gin.H {
	out := make([]gin.H, 0, len(list))
	for _, item := range list {
		out = append(out, toContentTemplateResponse(item))
	}
	return out
}

func toContentTemplateResponse(item domain.ContentTemplate) gin.H {
	return gin.H{
		"id":         item.ID,
		"name":       item.Name,
		"kind":       item.Kind,
		"status":     item.Status,
		"content":    decodeAnyJSON(item.ContentJSON),
		"created_at": item.CreatedAt,
		"updated_at": item.UpdatedAt,
	}
}
