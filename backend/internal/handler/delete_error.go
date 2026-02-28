package handler

import (
	"strings"

	"net/http"

	"github.com/cruisebooking/backend/internal/pkg/errcode"
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func isDeleteConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "foreign key constraint") ||
		strings.Contains(msg, "violates foreign key") ||
		strings.Contains(msg, "constraint failed") ||
		strings.Contains(msg, "still referenced")
}

func respondDeleteError(c *gin.Context, err error, resource string) {
	if isDeleteConstraintError(err) {
		name := strings.TrimSpace(resource)
		if name == "" {
			name = "resource"
		}
		response.Error(c, http.StatusConflict, errcode.ErrConflict, name+" has dependent records, cannot delete")
		return
	}
	response.InternalError(c, err)
}
