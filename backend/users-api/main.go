package main

import (
	"log"
	"os"

	"users-api/config"
	"users-api/handlers"
	"users-api/middleware"

	"github.com/gin-gonic/gin"
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

	// Servir archivos estáticos (avatares)
	router.Static("/uploads", "./uploads")

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
				profile.GET("", handlers.GetProfile)                  // GET /api/v1/profile
				profile.PUT("", handlers.UpdateProfile)               // PUT /api/v1/profile
				profile.PUT("/password", handlers.ChangePassword)     // PUT /api/v1/profile/password
				profile.POST("/avatar", handlers.UploadAvatar)        // POST /api/v1/profile/avatar
				profile.DELETE("/avatar", handlers.DeleteAvatar)      // DELETE /api/v1/profile/avatar
			}
		}
	}

	// Ruta raíz
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Users API - Sports Activities Platform",
			"version": "2.0.0",
			"features": []string{
				"JWT Authentication",
				"Extended User Profiles",
				"Avatar Upload/Management",
				"Role-based Access (Coming Soon)",
				"Email Verification (Coming Soon)",
			},
			"endpoints": gin.H{
				// Public endpoints
				"health":         "GET /api/v1/health",
				"auth_register":  "POST /api/v1/auth/register",
				"auth_login":     "POST /api/v1/auth/login",
				"auth_refresh":   "POST /api/v1/auth/refresh",
				"users_list":     "GET /api/v1/users",
				"user_by_id":     "GET /api/v1/users/:id",
				
				// Protected profile endpoints
				"profile":            "GET /api/v1/profile (protected)",
				"update_profile":     "PUT /api/v1/profile (protected)",
				"change_password":    "PUT /api/v1/profile/password (protected)",
				"upload_avatar":      "POST /api/v1/profile/avatar (protected)",
				"delete_avatar":      "DELETE /api/v1/profile/avatar (protected)",
				
				// Static files
				"avatars":            "GET /uploads/avatars/:filename",
			},
			"profile_fields": []string{
				"avatar_url", "bio", "phone", "birth_date", "location", "gender",
				"height", "weight", "sports_interests", "fitness_level", "social_links",
			},
		})
	})

	// Obtener puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 Users API v2.0.0 starting on port %s", port)
	log.Printf("📊 Database connected and migrated successfully")
	log.Printf("🔐 JWT authentication enabled")
	log.Printf("👤 Extended user profiles with avatar support")
	log.Printf("📸 Avatar upload/management enabled")
	log.Printf("📋 Available endpoints:")
	log.Printf("   • API Documentation: http://localhost:%s/", port)
	log.Printf("   • Health Check: http://localhost:%s/api/v1/health", port)
	log.Printf("   • Register: POST http://localhost:%s/api/v1/auth/register", port)
	log.Printf("   • Login: POST http://localhost:%s/api/v1/auth/login", port)
	log.Printf("   • Profile: GET/PUT http://localhost:%s/api/v1/profile", port)
	log.Printf("   • Avatar Upload: POST http://localhost:%s/api/v1/profile/avatar", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
