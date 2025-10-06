package middleware

import (
	"net/http"
	"strings"

	"users-api/config"
	"users-api/models"

	"github.com/gin-gonic/gin"
)

// RequireRole middleware que requiere un rol específico
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener información del usuario del contexto (viene del JWTAuth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "User not authenticated",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		// Obtener el usuario de la base de datos para verificar su rol
		var user models.User
		db := config.GetDB()
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "User not found",
				"message": "User account not found",
			})
			c.Abort()
			return
		}

		// Verificar si el usuario tiene el rol requerido
		if !hasRole(user.Role, requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient permissions",
				"message": "You don't have permission to access this resource",
				"required_role": requiredRole,
				"user_role": user.Role,
			})
			c.Abort()
			return
		}

		// Agregar información del rol al contexto
		c.Set("userRole", user.Role)
		c.Next()
	}
}

// RequireAnyRole middleware que requiere cualquiera de los roles especificados
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "User not authenticated",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}

		var user models.User
		db := config.GetDB()
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "User not found",
				"message": "User account not found",
			})
			c.Abort()
			return
		}

		// Verificar si el usuario tiene alguno de los roles permitidos
		hasPermission := false
		for _, role := range roles {
			if hasRole(user.Role, role) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient permissions",
				"message": "You don't have permission to access this resource",
				"allowed_roles": roles,
				"user_role": user.Role,
			})
			c.Abort()
			return
		}

		c.Set("userRole", user.Role)
		c.Next()
	}
}

// hasRole verifica si un usuario tiene un rol específico
// Implementa una jerarquía simple: super_admin > admin > moderator > user
func hasRole(userRole, requiredRole string) bool {
	// Normalizar roles a minúsculas
	userRole = strings.ToLower(strings.TrimSpace(userRole))
	requiredRole = strings.ToLower(strings.TrimSpace(requiredRole))

	// Si es exactamente el mismo rol
	if userRole == requiredRole {
		return true
	}

	// Jerarquía de roles (roles superiores incluyen permisos de roles inferiores)
	roleHierarchy := map[string]int{
		"user":        1,
		"moderator":   2,
		"admin":       3,
		"super_admin": 4,
		"root":        5,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	// Si ambos roles existen en la jerarquía, verificar nivel
	if userExists && requiredExists {
		return userLevel >= requiredLevel
	}

	// Si no están en la jerarquía, solo permitir coincidencia exacta
	return userRole == requiredRole
}

// IsAdmin helper function para verificar si un usuario es admin
func IsAdmin(userRole string) bool {
	return hasRole(userRole, "admin")
}

// IsRoot helper function para verificar si un usuario es root
func IsRoot(userRole string) bool {
	return hasRole(userRole, "root")
}

// GetUserRole obtiene el rol del usuario del contexto
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("userRole")
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	if !ok {
		return "", false
	}
	return roleStr, true
}
