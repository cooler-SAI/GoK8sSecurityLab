package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds JWT configuration parameters
type JWTConfig struct {
	SecretKey        string
	Issuer           string
	Audience         string
	TokenExpiry      time.Duration
	ValidateAudience bool
	ValidateIssuer   bool
}

// Claims structure for JWT claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// Context keys for storing authentication data
type contextKey string

const (
	UserIDKey          contextKey = "user_id"
	UsernameKey        contextKey = "username"
	RolesKey           contextKey = "roles"
	ClaimsKey          contextKey = "claims"
	IsAuthenticatedKey contextKey = "is_authenticated"
)

// DefaultJWTConfig returns default JWT configuration
func DefaultJWTConfig() JWTConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key-change-this-in-production"
		log.Println("⚠️ WARNING: Using default JWT secret. Set JWT_SECRET environment variable!")
	}

	return JWTConfig{
		SecretKey:        secret,
		Issuer:           "secure-halloween-server",
		Audience:         "halloween-api-users",
		TokenExpiry:      15 * time.Minute,
		ValidateAudience: true,
		ValidateIssuer:   true,
	}
}

// JWTAuth middleware for JWT token validation
func JWTAuth(config JWTConfig, optional bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				if optional {
					// If optional=true, continue without authentication
					ctx := context.WithValue(r.Context(), IsAuthenticatedKey, false)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				sendJSONError(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check "Bearer <token>" format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				sendJSONError(w, "Invalid authorization format. Use: Bearer <token>", http.StatusBadRequest)
				return
			}

			tokenString := parts[1]

			// Parse and validate token
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Validate signing algorithm
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.SecretKey), nil
			})

			if err != nil {
				log.Printf("JWT validation error: %v", err)
				sendJSONError(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				sendJSONError(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Additional claims validation
			if config.ValidateIssuer && claims.Issuer != config.Issuer {
				sendJSONError(w, "Invalid issuer", http.StatusUnauthorized)
				return
			}

			if config.ValidateAudience {
				audience, err := claims.GetAudience()
				if err != nil || !contains(audience, config.Audience) {
					sendJSONError(w, "Invalid audience", http.StatusUnauthorized)
					return
				}
			}

			// Check expiration time
			if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
				sendJSONError(w, "Token expired", http.StatusUnauthorized)
				return
			}

			// Add claims to context
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)
			ctx = context.WithValue(ctx, RolesKey, claims.Roles)
			ctx = context.WithValue(ctx, ClaimsKey, claims)
			ctx = context.WithValue(ctx, IsAuthenticatedKey, true)

			// Update request with new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GenerateToken creates a new JWT token
func GenerateToken(config JWTConfig, userID, username string, roles []string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.Issuer,
			Audience:  jwt.ClaimStrings{config.Audience},
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// GetUserID Helper functions for extracting data from context
func GetUserID(ctx context.Context) (string, bool) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return "", false
	}
	return val.(string), true
}

func GetUsername(ctx context.Context) (string, bool) {
	val := ctx.Value(UsernameKey)
	if val == nil {
		return "", false
	}
	return val.(string), true
}

func GetRoles(ctx context.Context) ([]string, bool) {
	val := ctx.Value(RolesKey)
	if val == nil {
		return nil, false
	}
	return val.([]string), true
}

func IsAuthenticated(ctx context.Context) bool {
	val := ctx.Value(IsAuthenticatedKey)
	if val == nil {
		return false
	}
	return val.(bool)
}

// RequireRole middleware checks for specific role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := GetRoles(r.Context())
			if !ok {
				sendJSONError(w, "No roles found", http.StatusForbidden)
				return
			}

			if !contains(roles, role) {
				sendJSONError(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Utility functions
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := fmt.Fprintf(w, `{"error": "%s"}`, message)
	if err != nil {
		return
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
