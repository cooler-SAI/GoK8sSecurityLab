package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"websecure/handler"
	"websecure/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// JWT configuration
	jwtConfig := middleware.DefaultJWTConfig()

	// Middleware chains for different endpoints

	// Public endpoints (no authentication required)
	healthChain := handler.RateLimit(
		middleware.SecurityHeaders(
			http.HandlerFunc(handler.HealthHandler),
		),
	)

	halloweenChain := handler.RateLimit(
		middleware.SecurityHeaders(
			http.HandlerFunc(handler.HalloweenHandler),
		),
	)

	apiChain := handler.RateLimit(
		middleware.SecurityHeaders(
			http.HandlerFunc(handler.HalloweenAPIHandler),
		),
	)

	infoChain := handler.RateLimit(
		middleware.SecurityHeaders(
			http.HandlerFunc(handler.InfoHandler),
		),
	)

	// Authentication endpoints
	loginChain := handler.RateLimit(
		middleware.SecurityHeaders(
			http.HandlerFunc(handler.JWTHandler),
		),
	)

	// Protected endpoints (require JWT)
	protectedChain := handler.RateLimit(
		middleware.SecurityHeaders(
			middleware.JWTAuth(jwtConfig, false)(
				http.HandlerFunc(handler.ProtectedHandler),
			),
		),
	)

	// Admin endpoint (requires admin role)
	adminChain := handler.RateLimit(
		middleware.SecurityHeaders(
			middleware.JWTAuth(jwtConfig, false)(
				middleware.RequireRole("admin")(
					http.HandlerFunc(handler.AdminHandler),
				),
			),
		),
	)

	// Endpoint with optional authentication
	optionalAuthChain := handler.RateLimit(
		middleware.SecurityHeaders(
			middleware.JWTAuth(jwtConfig, true)(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if middleware.IsAuthenticated(r.Context()) {
						handler.ProtectedHandler(w, r)
					} else {
						handler.HalloweenAPIHandler(w, r)
					}
				}),
			),
		),
	)

	// Secure greet endpoint
	secureGreetChain := handler.RateLimit(
		middleware.SecurityHeaders(
			http.HandlerFunc(handler.SecureGreetHandler),
		),
	)

	// Vulnerable greet endpoint (for demonstration)
	vulnerableGreetChain := handler.RateLimit(
		middleware.SecurityHeaders(
			http.HandlerFunc(handler.VulnerableGreetHandler),
		),
	)

	// Register routes
	http.Handle("/", healthChain)
	http.Handle("/halloween", halloweenChain)
	http.Handle("/api/halloween", apiChain)
	http.Handle("/info", infoChain)
	http.Handle("/greet", secureGreetChain)
	http.Handle("/vulnerable-greet", vulnerableGreetChain)

	// JWT Authentication routes
	http.Handle("/api/auth/jwt", loginChain)        // POST for login, GET for token info
	http.Handle("/api/protected", protectedChain)   // Requires JWT
	http.Handle("/api/admin", adminChain)           // Requires admin role
	http.Handle("/api/optional", optionalAuthChain) // Optional authentication

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Printf("🎃 Starting Secure Halloween Server on port %s", port)
	log.Printf("\n📋 AVAILABLE ENDPOINTS:")
	log.Printf("\n🔓 PUBLIC ENDPOINTS:")
	log.Printf("   🌐 Health Check:        http://localhost:%s/", port)
	log.Printf("   🎃 Halloween Page:      http://localhost:%s/halloween", port)
	log.Printf("   🔗 Halloween API:       http://localhost:%s/api/halloween", port)
	log.Printf("   ℹ️  Server Info:         http://localhost:%s/info", port)
	log.Printf("   ✅ SECURE GREET:         http://localhost:%s/greet?name=YourName", port)
	log.Printf("   ⚠️  VULNERABLE GREET:    http://localhost:%s/vulnerable-greet?name=YourName", port)

	log.Printf("\n🔐 JWT AUTHENTICATION:")
	log.Printf("   🔑 Login (POST):        http://localhost:%s/api/auth/jwt", port)
	log.Printf("   📋 Token Info (GET):    http://localhost:%s/api/auth/jwt", port)
	log.Printf("   🛡️  Protected API:       http://localhost:%s/api/protected", port)
	log.Printf("   👑 Admin API:           http://localhost:%s/api/admin", port)
	log.Printf("   🤔 Optional Auth:       http://localhost:%s/api/optional", port)

	log.Printf("\n🔧 JWT CONFIGURATION:")
	log.Printf("   Issuer: %s", jwtConfig.Issuer)
	log.Printf("   Audience: %s", jwtConfig.Audience)
	log.Printf("   Token Expiry: %v", jwtConfig.TokenExpiry)

	log.Printf("\n⚡ SECURITY FEATURES:")
	log.Printf("   Rate Limiting: 10 requests/second")
	log.Printf("   Security Headers: enabled")
	log.Printf("   JWT Authentication: enabled")
	log.Printf("   XSS Protection: enabled")

	log.Printf("\n👤 DEMO USERS:")
	log.Printf("   alice:password123  (roles: user, premium)")
	log.Printf("   bob:secure456      (roles: user)")
	log.Printf("   admin:admin789     (roles: admin, user)")

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error starting server: %v", err)
	}
}
