package middleware

import (
	"activities-api/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminOnly middleware verifica que el usuario tenga rol de admin
func AdminOnly(usersAPI string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Verificar el rol directamente del JWT (más eficiente, sin llamada a API)
		if claims.Role != "admin" && claims.Role != "root" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin role required"})
			c.Abort()
			return
		}

		// Setear información del usuario en el contexto de Gin y en el request
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		reqCtx := context.WithValue(c.Request.Context(), utils.ContextUserIDKey, claims.UserID)
		reqCtx = context.WithValue(reqCtx, utils.ContextUserRoleKey, claims.Role)
		c.Request = c.Request.WithContext(reqCtx)

		c.Next()
	}
}
