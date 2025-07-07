package handler

import (
	"net/http"
	"orders/internal/service"
	"orders/models"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	uid := c.Param("id")
	order, err := h.service.GetOrderByUID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}
func InitOrderRouter(r *gin.Engine, service *service.OrderService) {
	orderHandler := NewOrderHandler(service)
	r.GET("/order/:id", orderHandler.GetOrder)
	r.POST("/order", orderHandler.CreateOrder)
	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})
}
