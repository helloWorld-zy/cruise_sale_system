package handler

import (
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type RouteService interface{ List() ([]interface{}, error) }

type RouteHandler struct{ svc RouteService }

func NewRouteHandler(svc RouteService) *RouteHandler { return &RouteHandler{svc: svc} }

func (h *RouteHandler) List(c *gin.Context) {
	var list []interface{}
	if h.svc != nil {
		list, _ = h.svc.List()
	} else {
		list = []interface{}{}
	}
	response.Success(c, list)
}
