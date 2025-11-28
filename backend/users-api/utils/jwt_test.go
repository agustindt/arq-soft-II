package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	// Setup
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Test data
	userID := uint(1)
	email := "test@example.com"
	username := "testuser"
	role := "user"

	// Generate token
	token, err := GenerateJWT(userID, email, username, role)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT_ValidToken(t *testing.T) {
	// Setup
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Generate a valid token
	userID := uint(1)
	email := "test@example.com"
	username := "testuser"
	role := "admin"

	token, err := GenerateJWT(userID, email, username, role)
	assert.NoError(t, err)

	// Validate token
	claims, err := ValidateJWT(token)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	// Setup
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Create an expired token manually
	claims := Claims{
		UserID:   1,
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "sports-activities-api",
			Subject:   "user-authentication",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTSecret())
	assert.NoError(t, err)

	// Try to validate expired token
	validatedClaims, err := ValidateJWT(tokenString)

	// Assertions - should fail
	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
	// Setup
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Generate token with one secret
	token, err := GenerateJWT(1, "test@example.com", "testuser", "user")
	assert.NoError(t, err)

	// Change the secret
	os.Setenv("JWT_SECRET", "different-secret-key")

	// Try to validate with different secret
	claims, err := ValidateJWT(token)

	// Assertions - should fail
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateJWT_MalformedToken(t *testing.T) {
	// Setup
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Try to validate a malformed token
	claims, err := ValidateJWT("this.is.not.a.valid.token")

	// Assertions - should fail
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestRefreshJWT(t *testing.T) {
	// Setup
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Generate original token
	userID := uint(1)
	email := "test@example.com"
	username := "testuser"
	role := "user"

	originalToken, err := GenerateJWT(userID, email, username, role)
	assert.NoError(t, err)

	// Wait to ensure different timestamps (JWT precision is 1 second)
	time.Sleep(1100 * time.Millisecond)

	// Refresh token
	newToken, err := RefreshJWT(originalToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)
	
	// Note: Tokens might be the same if generated in the same second
	// The important thing is that refresh works and returns a valid token
	assert.NotEqual(t, originalToken, newToken, "Token should be different after waiting > 1 second")

	// Validate new token has same user data
	claims, err := ValidateJWT(newToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, role, claims.Role)
}

func TestRefreshJWT_ExpiredToken(t *testing.T) {
	// Setup
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Create an expired token
	claims := Claims{
		UserID:   1,
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "sports-activities-api",
			Subject:   "user-authentication",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTSecret())
	assert.NoError(t, err)

	// Try to refresh expired token
	newToken, err := RefreshJWT(tokenString)

	// Assertions - should fail
	assert.Error(t, err)
	assert.Empty(t, newToken)
}

func TestGetJWTSecret_WithEnv(t *testing.T) {
	expectedSecret := "my-test-secret"
	os.Setenv("JWT_SECRET", expectedSecret)
	defer os.Unsetenv("JWT_SECRET")

	secret := GetJWTSecret()
	assert.Equal(t, []byte(expectedSecret), secret)
}

func TestGetJWTSecret_WithoutEnv(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	secret := GetJWTSecret()
	assert.NotNil(t, secret)
	assert.Equal(t, []byte("your-super-secret-jwt-key-here"), secret)
}

// Integration test: Generate token with short expiration and verify it expires
func TestTokenExpiration_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// Create a token that expires in 2 seconds
	claims := Claims{
		UserID:   1,
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "sports-activities-api",
			Subject:   "user-authentication",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTSecret())
	assert.NoError(t, err)

	// Validate immediately - should work
	validatedClaims, err := ValidateJWT(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, validatedClaims)

	// Wait for token to expire
	t.Log("Waiting 3 seconds for token to expire...")
	time.Sleep(3 * time.Second)

	// Validate again - should fail
	validatedClaims, err = ValidateJWT(tokenString)
	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
	assert.Contains(t, err.Error(), "token is expired")
}

