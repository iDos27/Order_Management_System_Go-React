package main

import (
	"auth-service/internal/database"
	"auth-service/internal/handlers"
	"auth-service/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Ustawienie trybu produkcyjnego dla Gin
	gin.SetMode(gin.ReleaseMode)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5433/auth_service?sslmode=disable"
		log.Println("Using default DATABASE_URL")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret-key"
		log.Println("Using default JWT_SECRET")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		log.Println("Using default PORT 8081")
	}

	os.Setenv("AUTH_DATABASE_URL", dbURL)
	os.Setenv("JWT_SECRET", jwtSecret)

	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to the database")

	if err := db.CreateUsersTable(); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
	log.Println("Users table is ready")

	authHandler := handlers.NewAuthHandler(db)
	authMiddleware := middleware.NewAuthMiddleware(db)

	// Konfiguracja routera
	router := gin.Default()

	// Middleware CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// auth_request
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	// Routing API
	api := router.Group("/api/v1")
	{
		// Endpointy publiczne
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		api.GET("/validate-token", authMiddleware.ValidateTokenEndpoint)

		api.GET("/validate", func(c *gin.Context) {
			middleware := authMiddleware.RequireAuth()
			middleware(c)

			if !c.IsAborted() {
				c.Status(200) // NGINX auth_request success
			}
		})

		// Endpointy chronione
		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.GET("/profile", authHandler.Profile)
		}

		admin := api.Group("/admin")
		admin.Use(authMiddleware.RequireAuth())
		admin.Use(authMiddleware.RequireAdmin())
		{
			admin.GET("/users", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Admin access to users",
				})
			})
		}
	}

	// Uruchomienie serwera
	log.Printf("Starting Auth Service on port %s", port)
	log.Fatal(router.Run(":" + port))
}
