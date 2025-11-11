package middleware

import (
	"arq-soft-II/backend/reservations-api/clients"
	"net/http"
	"strings"

	"arq-soft-II/backend/reservations-api/utils"

	"github.com/gin-gonic/gin"
)

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

		user, err := clients.GetUserByID(usersAPI, claims.UserID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not accessible"})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin role required"})
			c.Abort()
			return
		}

		// Setear el user_id en el contexto si lo querés usar después
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
