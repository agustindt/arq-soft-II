package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"arq-soft-II/config/httpx"

	"github.com/golang-jwt/jwt/v5"
)

// Función para generar un JWT de prueba
func createTestJWT(secret string, isAdmin bool) (string, error) {
	claims := httpx.Claims{
		UserID:  "12345",
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func main() {
	os.Setenv("JWT_SECRET", "your-super-secret-jwt-key-here")

	// Generar un JWT de prueba
	token, err := createTestJWT(os.Getenv("JWT_SECRET"), true)
	if err != nil {
		panic(err)
	}
	fmt.Println("Token JWT generado:")
	fmt.Println(token)

	// Simular un request HTTP con el token
	req := httptest.NewRequest(http.MethodGet, "/admin/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	// Crear un handler de prueba (solo admin)
	handler := httpx.RequireAuth(httpx.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Acceso autorizado a endpoint admin")
	})))

	handler.ServeHTTP(rr, req)

	fmt.Println("Respuesta del middleware:", rr.Body.String())
}
