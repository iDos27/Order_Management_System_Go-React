package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/iDos27/order-management/order-service/internal/database"
	"github.com/iDos27/order-management/order-service/internal/handlers"
	"github.com/iDos27/order-management/order-service/internal/websocket"
)

func main() {
	// Połączenie z bazą
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Inicjalizacja WebSocket hub
	hub := websocket.NewHub()
	go hub.Run() // Uruchamiamy hub w goroutine

	// Inicjalizacja handlers
	orderHandler := handlers.NewOrderHandler(db, hub)

	// Setup routera
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost", "http://localhost:80"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "order-service"})
	})

	// Routers
	api := router.Group("/api")
	{
		api.GET("/orders", orderHandler.GetAllOrders)
		api.GET("/orders/:id", orderHandler.GetOrderByID)
		api.POST("/orders", orderHandler.CreateOrder)
		api.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
	}

	// WebSocket endpoint
	router.GET("/ws", websocket.HandleWebSocket(hub))

	// Uruchomienie serwera
	port := getEnv("SERVER_PORT", "8080")
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
