package handler

import "github.com/gin-gonic/gin"

type CabinHandler struct{}

func NewCabinHandler() *CabinHandler        { return &CabinHandler{} }
func (h *CabinHandler) List(c *gin.Context) { c.JSON(200, gin.H{"data": []interface{}{}}) }
