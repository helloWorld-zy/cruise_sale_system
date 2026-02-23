package handler

import (
	"github.com/cruisebooking/backend/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler { return &UserHandler{} }

func (h *UserHandler) Login(c *gin.Context) {
	response.Success(c, gin.H{"token": "stub"})
}

func (h *UserHandler) Profile(c *gin.Context) {
	response.Success(c, gin.H{"id": 1})
}
