package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cruisebooking/backend/internal/domain"
	"github.com/cruisebooking/backend/internal/middleware"
	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type StaffService interface {
	Create(ctx context.Context, name, email, role string) (*domain.Staff, error)
	AssignRole(ctx context.Context, id int64, role string, operatorID int64) error
	List(ctx context.Context) ([]domain.Staff, error)
	GetByID(ctx context.Context, id int64) (*domain.Staff, error)
	Update(ctx context.Context, staff *domain.Staff) error
	Delete(ctx context.Context, id int64) error
}

type StaffHandler struct {
	svc StaffService
}

func NewStaffHandler(svc StaffService) *StaffHandler {
	return &StaffHandler{svc: svc}
}

func (h *StaffHandler) Create(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Role  string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	staff, err := h.svc.Create(c.Request.Context(), req.Name, req.Email, req.Role)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	response.Success(c, staff)
}

func (h *StaffHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, list)
}

func (h *StaffHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	staff, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if staff == nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "staff not found")
		return
	}
	response.Success(c, staff)
}

func (h *StaffHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	staff, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if staff == nil {
		response.Error(c, http.StatusNotFound, errcode.ErrNotFound, "staff not found")
		return
	}

	if req.Name != "" {
		staff.RealName = req.Name
	}
	if req.Email != "" {
		staff.Email = req.Email
	}
	if req.Role != "" {
		staff.Role = req.Role
	}

	if err := h.svc.Update(c.Request.Context(), staff); err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, staff)
}

func (h *StaffHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.InternalError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *StaffHandler) AssignRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, "invalid id")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}

	if err := h.svc.AssignRole(c.Request.Context(), id, req.Role, parseStaffID(c)); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.ErrValidation, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id, "role": req.Role})
}

func parseStaffID(c *gin.Context) int64 {
	v, ok := c.Get(middleware.ContextKeyStaffID)
	if !ok {
		return 0
	}
	switch value := v.(type) {
	case string:
		if id, err := strconv.ParseInt(value, 10, 64); err == nil {
			return id
		}
	case int64:
		return value
	}
	return 0
}
