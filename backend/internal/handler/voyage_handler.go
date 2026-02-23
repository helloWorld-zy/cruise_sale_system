package handler

import "github.com/gin-gonic/gin"

type VoyageHandler struct{}

func NewVoyageHandler() *VoyageHandler       { return &VoyageHandler{} }
func (h *VoyageHandler) List(c *gin.Context) { c.JSON(200, gin.H{"data": []interface{}{}}) }
