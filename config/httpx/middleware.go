package httpx

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Claims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// Definir tipo de clave único para el contexto
type contextKey string

const userContextKey contextKey = "user"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			writeError(w, "Falta token de autorización", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			writeError(w, "Error interno: falta JWT_SECRET", http.StatusInternalServerError)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			writeError(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userContextKey).(*Claims)
		if !ok || !claims.IsAdmin {
			writeError(w, "Acceso restringido a administradores", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireServiceToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceToken := r.Header.Get("X-Service-Token")
		expected := os.Getenv("SERVICE_TOKEN")

		if expected == "" {
			writeError(w, "Error interno: falta SERVICE_TOKEN", http.StatusInternalServerError)
			return
		}

		if serviceToken != expected {
			writeError(w, "Token interno inválido", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ErrorResponse{
		Status:  "error",
		Message: message,
	}
	json.NewEncoder(w).Encode(resp)
	log.Printf("%s (%d)", message, code)
}
