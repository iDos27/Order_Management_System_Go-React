package handlers

import (
	"net/http"
	"strconv"

	"order-management-system/internal/database"
	"order-management-system/internal/models"
	"order-management-system/internal/websocket"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	db  *database.DB
	hub *websocket.Hub
}

func NewOrderHandler(db *database.DB, hub *websocket.Hub) *OrderHandler {
	return &OrderHandler{db: db, hub: hub}
}

// GET /api/orders - Lista wszystkich zamówień (admin)
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	rows, err := h.db.Query(`
        SELECT id, customer_name, customer_email, status, total_amount, created_at, updated_at 
        FROM orders 
        ORDER BY created_at DESC
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.CustomerName, &order.CustomerEmail,
			&order.Status, &order.TotalAmount, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order"})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

// GET /api/orders/:id - Pojedycze zamówienie (klient)
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	err = h.db.QueryRow(`
        SELECT id, customer_name, customer_email, status, total_amount, created_at, updated_at 
        FROM orders WHERE id = $1
    `, id).Scan(&order.ID, &order.CustomerName, &order.CustomerEmail,
		&order.Status, &order.TotalAmount, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// POST /api/orders - Tworzenie nowego zamówienia
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err := h.db.QueryRow(`
		INSERT INTO orders (customer_name, customer_email, status, total_amount)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
		`, order.CustomerName, order.CustomerEmail, models.StatusNew, order.TotalAmount).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// PATCH /api/orders/:id/status - Zmiana statusu zamówienia
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Walidacja
	validStatuses := []string{"new", "confirmed", "shipped", "delivered", "cancelled"}
	isValid := false
	for _, status := range validStatuses {
		if statusUpdate.Status == status {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	// Aktualizacja statusu w bazie
	_, err = h.db.Exec(`
        UPDATE orders 
        SET status = $1, updated_at = CURRENT_TIMESTAMP 
        WHERE id = $2
    `, statusUpdate.Status, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// Powiadomienie przez WebSocket
	h.hub.BroadcastOrderUpdate(id, statusUpdate.Status, "admin")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Order status updated successfully",
		"order_id":   id,
		"new_status": statusUpdate.Status,
	})
}
