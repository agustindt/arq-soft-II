package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"users-api/config"
	"users-api/handlers"
	"users-api/middleware"
)

func main() {
	// Conectar a la base de datos
	config.ConnectDatabase()

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware CORS básico
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

	// Rutas públicas
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Users API is running",
				"service": "users-api",
			})
		})

		// Rutas de autenticación (públicas)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		// Rutas públicas de usuarios
		users := api.Group("/users")
		{
			users.GET("", handlers.ListUsers)           // GET /api/v1/users
			users.GET("/:id", handlers.GetUserByID)     // GET /api/v1/users/:id
		}

		// Rutas protegidas (requieren JWT)
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			// Perfil del usuario autenticado
			profile := protected.Group("/profile")
			{
				profile.GET("", handlers.GetProfile)          // GET /api/v1/profile
				profile.PUT("", handlers.UpdateProfile)       // PUT /api/v1/profile
				profile.PUT("/password", handlers.ChangePassword) // PUT /api/v1/profile/password
			}
		}
	}

	// Ruta raíz
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Users API - Sports Activities Platform",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":         "GET /api/v1/health",
				"auth_register":  "POST /api/v1/auth/register",
				"auth_login":     "POST /api/v1/auth/login",
				"auth_refresh":   "POST /api/v1/auth/refresh",
				"users_list":     "GET /api/v1/users",
				"user_by_id":     "GET /api/v1/users/:id",
				"profile":        "GET /api/v1/profile (protected)",
				"update_profile": "PUT /api/v1/profile (protected)",
				"change_password": "PUT /api/v1/profile/password (protected)",
			},
		})
	})

	// Obtener puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 Users API starting on port %s", port)
	log.Printf("📊 Database connected and migrated successfully")
	log.Printf("🔐 JWT authentication enabled")
	log.Printf("📋 Available endpoints:")
	log.Printf("   • Health: http://localhost:%s/api/v1/health", port)
	log.Printf("   • Register: POST http://localhost:%s/api/v1/auth/register", port)
	log.Printf("   • Login: POST http://localhost:%s/api/v1/auth/login", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
