// Package main implements the Users API microservice.
//
// The Users API is responsible for user authentication, authorization, and profile management.
// It provides JWT-based authentication with role-based access control (RBAC) supporting
// hierarchical roles: user, moderator, admin, super_admin, and root.
//
// Key Features:
//   - User registration and login with bcrypt password hashing
//   - JWT token generation with role claims for efficient authorization
//   - Extended user profiles with sports-specific fields
//   - Avatar upload and management
//   - Public and private profile data separation
//   - Admin user management (create, update roles/status, view stats)
//   - Root user capabilities (user deletion)
//
// Database: MySQL 8.0 with GORM
// Port: 8081
// API Version: 2.1.0
//
// For complete API documentation, see docs/api/users-api.md
package main

import (
	"log"
	"os"

	"users-api/config"
	"users-api/controllers"
	"users-api/middleware"
	"users-api/repositories"
	"users-api/services"

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

	// Inicializar dependencias
	userRepo := repositories.NewUserRepository(config.GetDB())
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	adminService := services.NewAdminService(userRepo)

	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	adminController := controllers.NewAdminController(adminService)

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
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
		}

		// Ruta especial para crear usuario root (solo disponible si no existe root)
		api.POST("/admin/create-root", adminController.CreateRoot)

		// Rutas públicas de usuarios
		users := api.Group("/users")
		{
			users.GET("", userController.ListUsers)       // GET /api/v1/users
			users.GET("/:id", userController.GetUserByID) // GET /api/v1/users/:id
		}

		// Rutas protegidas (requieren JWT)
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			// Perfil del usuario autenticado
			profile := protected.Group("/profile")
			{
				profile.GET("", userController.GetProfile)              // GET /api/v1/profile
				profile.PUT("", userController.UpdateProfile)           // PUT /api/v1/profile
				profile.PUT("/password", userController.ChangePassword) // PUT /api/v1/profile/password
				profile.POST("/avatar", userController.UploadAvatar)    // POST /api/v1/profile/avatar
				profile.DELETE("/avatar", userController.DeleteAvatar)  // DELETE /api/v1/profile/avatar
			}
		}

		// Rutas de administración (requieren JWT + rol admin)
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth())
		admin.Use(middleware.RequireRole("admin"))
		{
			// Gestión de usuarios
			admin.GET("/users", adminController.ListAllUsers)                // GET /api/v1/admin/users
			admin.POST("/users", adminController.CreateUser)                 // POST /api/v1/admin/users
			admin.PUT("/users/:id/role", adminController.UpdateUserRole)     // PUT /api/v1/admin/users/:id/role
			admin.PUT("/users/:id/status", adminController.UpdateUserStatus) // PUT /api/v1/admin/users/:id/status
			admin.GET("/stats", adminController.GetSystemStats)              // GET /api/v1/admin/stats
		}

		// Rutas solo para root users
		root := api.Group("/admin")
		root.Use(middleware.JWTAuth())
		root.Use(middleware.RequireRole("root"))
		{
			root.DELETE("/users/:id", adminController.DeleteUser) // DELETE /api/v1/admin/users/:id
		}
	}

	// Ruta raíz
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Users API - Sports Activities Platform",
			"version": "2.1.0",
			"features": []string{
				"JWT Authentication",
				"Extended User Profiles",
				"Avatar Upload/Management",
				"Role-based Access Control",
				"Admin User Management",
				"Email Verification (Coming Soon)",
			},
			"endpoints": gin.H{
				// Public endpoints
				"health":        "GET /api/v1/health",
				"auth_register": "POST /api/v1/auth/register",
				"auth_login":    "POST /api/v1/auth/login",
				"auth_refresh":  "POST /api/v1/auth/refresh",
				"users_list":    "GET /api/v1/users",
				"user_by_id":    "GET /api/v1/users/:id",

				// Protected profile endpoints
				"profile":         "GET /api/v1/profile (protected)",
				"update_profile":  "PUT /api/v1/profile (protected)",
				"change_password": "PUT /api/v1/profile/password (protected)",
				"upload_avatar":   "POST /api/v1/profile/avatar (protected)",
				"delete_avatar":   "DELETE /api/v1/profile/avatar (protected)",

				// Admin endpoints (admin role required)
				"create_root":         "POST /api/v1/admin/create-root (public, secret key required)",
				"admin_users_list":    "GET /api/v1/admin/users (admin)",
				"admin_create_user":   "POST /api/v1/admin/users (admin)",
				"admin_update_role":   "PUT /api/v1/admin/users/:id/role (admin)",
				"admin_update_status": "PUT /api/v1/admin/users/:id/status (admin)",
				"admin_stats":         "GET /api/v1/admin/stats (admin)",
				"admin_delete_user":   "DELETE /api/v1/admin/users/:id (root only)",

				// Static files
				"avatars": "GET /uploads/avatars/:filename",
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

	log.Printf("🚀 Users API v2.1.0 starting on port %s", port)
	log.Printf("📊 Database connected and migrated successfully")
	log.Printf("🔐 JWT authentication enabled")
	log.Printf("👤 Extended user profiles with avatar support")
	log.Printf("📸 Avatar upload/management enabled")
	log.Printf("🛡 Role-based access control (user/moderator/admin/root)")
	log.Printf("👨‍💼 Admin user management system")
	log.Printf("📋 Available endpoints:")
	log.Printf("   • API Documentation: http://localhost:%s/", port)
	log.Printf("   • Health Check: http://localhost:%s/api/v1/health", port)
	log.Printf("   • Register: POST http://localhost:%s/api/v1/auth/register", port)
	log.Printf("   • Login: POST http://localhost:%s/api/v1/auth/login", port)
	log.Printf("   • Profile: GET/PUT http://localhost:%s/api/v1/profile", port)
	log.Printf("   • Avatar Upload: POST http://localhost:%s/api/v1/profile/avatar", port)
	log.Printf("   • Create Root: POST http://localhost:%s/api/v1/admin/create-root", port)
	log.Printf("   • Admin Panel: GET http://localhost:%s/api/v1/admin/users", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
