package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// HalloweenHandler - holiday page
func HalloweenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html>
<head>
    <title>🎃 Secure Halloween Server</title>
    <style>
        body { 
            background: #1a1a2e; 
            color: #ff7b25; 
            font-family: 'Courier New', monospace;
            text-align: center;
            padding: 50px;
        }
        .pumpkin { 
            font-size: 80px; 
            margin: 20px;
            text-shadow: 0 0 10px #ff4500;
            animation: glow 2s infinite alternate;
        }
        @keyframes glow {
            from { text-shadow: 0 0 10px #ff4500; }
            to { text-shadow: 0 0 20px #ff0000, 0 0 30px #ff4500; }
        }
        .security-badge {
            background: #162447;
            padding: 20px;
            border-radius: 15px;
            border: 2px solid #ff7b25;
            margin: 30px auto;
            max-width: 500px;
            box-shadow: 0 0 15px rgba(255, 123, 37, 0.3);
        }
        .links a {
            display: inline-block;
            margin: 10px;
            padding: 10px 20px;
            background: #e94560;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            transition: background 0.3s;
        }
        .links a:hover {
            background: #ff7b25;
        }
    </style>
</head>
<body>
    <div class="pumpkin">🎃</div>
    <h1>HAPPY HALLOWEEN!</h1>
    
    <div class="security-badge">
        <h2>🔒 Secure Go Server</h2>
        <p><strong>Protected by:</strong></p>
        <p>✅ Rate Limiting</p>
        <p>✅ Security Headers</p>
        <p><small>Client IP: ` + r.RemoteAddr + `</small></p>
    </div>

    <div class="links">
        <a href="/">Health Check</a>
        <a href="/api/halloween">API</a>
        <a href="/info">Server Info</a>
        <a href="/greet?name=Test">Vulnerable Greet</a>
    </div>

    <div style="margin-top: 30px;">
        <p><em>No tricks, only treats! 🍬</em></p>
    </div>
</body>
</html>`

	_, err := fmt.Fprint(w, html)
	if err != nil {
		log.Printf("Error writing Halloween response: %v", err)
	}
}

// HalloweenAPIHandler - JSON API for Halloween
func HalloweenAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":    "spookily_secure",
		"message":   "Happy Halloween! 🎃",
		"client_ip": r.RemoteAddr,
		"security":  "protected_by_middleware",
		"features": []string{
			"rate_limiting",
			"security_headers",
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
