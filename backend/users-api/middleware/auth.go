package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"users-api/utils"
)

// JWTAuth middleware para validar JWT tokens
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header required",
				"message": "Please provide a valid JWT token",
			})
			c.Abort()
			return
		}

		// Verificar que el token tenga el formato correcto (Bearer <token>)
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header format",
				"message": "Authorization header must be 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validar el token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Agregar información del usuario al contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.Username)

		// Continuar con el siguiente handler
		c.Next()
	}
}

// GetUserFromContext obtiene la información del usuario desde el contexto
func GetUserFromContext(c *gin.Context) (userID uint, email, username string, exists bool) {
	userIDInterface, existsID := c.Get("user_id")
	emailInterface, existsEmail := c.Get("user_email")
	usernameInterface, existsUsername := c.Get("user_username")

	if !existsID || !existsEmail || !existsUsername {
		return 0, "", "", false
	}

	userID, okID := userIDInterface.(uint)
	email, okEmail := emailInterface.(string)
	username, okUsername := usernameInterface.(string)

	if !okID || !okEmail || !okUsername {
		return 0, "", "", false
	}

	return userID, email, username, true
}
