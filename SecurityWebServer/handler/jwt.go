package handler

import (
	"SecurityWebServer/middleware"
	"encoding/json"
	"net/http"
	"time"
)

// JWTResponse response with JWT token
type JWTResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
}

// LoginRequest login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DemoUser demo user
type DemoUser struct {
	ID       string
	Username string
	Password string
	Roles    []string
}

// Demo users for testing
var demoUsers = map[string]DemoUser{
	"alice": {
		ID:       "user_001",
		Username: "alice",
		Password: "password123", // In real applications, use bcrypt!
		Roles:    []string{"user", "premium"},
	},
	"bob": {
		ID:       "user_002",
		Username: "bob",
		Password: "secure456",
		Roles:    []string{"user"},
	},
	"admin": {
		ID:       "admin_001",
		Username: "admin",
		Password: "admin789",
		Roles:    []string{"admin", "user"},
	},
}

// JWTHandler handler for JWT operations
func JWTHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleLogin(w, r)
	case http.MethodGet:
		handleTokenInfo(w, r)
	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handleLogin processes user login and returns JWT token
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, ok := demoUsers[req.Username]
	if !ok || user.Password != req.Password {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	config := middleware.DefaultJWTConfig()
	token, err := middleware.GenerateToken(config, user.ID, user.Username, user.Roles)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Return response
	response := JWTResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int(config.TokenExpiry.Seconds()),
		UserID:    user.ID,
		Username:  user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// handleTokenInfo returns information about the current JWT token
func handleTokenInfo(w http.ResponseWriter, r *http.Request) {
	// Check authentication via middleware
	if !middleware.IsAuthenticated(r.Context()) {
		http.Error(w, `{"error": "Not authenticated"}`, http.StatusUnauthorized)
		return
	}

	// Extract information from context
	userID, _ := middleware.GetUserID(r.Context())
	username, _ := middleware.GetUsername(r.Context())
	roles, _ := middleware.GetRoles(r.Context())

	response := map[string]interface{}{
		"authenticated": true,
		"user_id":       userID,
		"username":      username,
		"roles":         roles,
		"timestamp":     time.Now().Format(time.RFC3339),
		"message":       "🎃 Your JWT token is valid!",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// ProtectedHandler protected endpoint handler (requires JWT)
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r.Context()) {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}

	userID, _ := middleware.GetUserID(r.Context())
	username, _ := middleware.GetUsername(r.Context())

	response := map[string]interface{}{
		"message":   "🎃 Welcome to the protected Halloween API!",
		"user_id":   userID,
		"username":  username,
		"endpoint":  "Protected Resource",
		"timestamp": time.Now().Format(time.RFC3339),
		"security":  "This endpoint requires valid JWT token",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// AdminHandler requires admin role
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	username, _ := middleware.GetUsername(r.Context())

	response := map[string]interface{}{
		"message":   "👑 Welcome, Administrator!",
		"username":  username,
		"endpoint":  "Admin Panel",
		"features":  []string{"User Management", "Security Logs", "System Configuration"},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
