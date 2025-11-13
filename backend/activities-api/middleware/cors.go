package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura los headers CORS para permitir peticiones del frontend
func CORSMiddleware(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if ctx.Request.Method == http.MethodOptions {
		ctx.Status(http.StatusNoContent)
		return
	}

	ctx.Next()
}
