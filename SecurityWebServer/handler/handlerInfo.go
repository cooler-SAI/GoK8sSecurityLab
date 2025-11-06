package handler

import (
	"fmt"
	"log"
	"net/http"
)

// InfoHandler - server information and endpoints (EXPORTED)
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	info := `<!DOCTYPE html>
<html>
<head>
    <title>🔒 Server Information</title>
    <style>
        body { 
            background: #0f3460; 
            color: #e94560; 
            font-family: Arial, sans-serif; 
            padding: 30px; 
            margin: 0;
        }
        .container { 
            max-width: 800px; 
            margin: 0 auto; 
            background: #16213e; 
            padding: 30px; 
            border-radius: 15px; 
            box-shadow: 0 0 20px rgba(0,0,0,0.3);
        }
        .endpoint { 
            background: #1a1a2e; 
            padding: 15px; 
            margin: 10px 0; 
            border-radius: 8px; 
            border-left: 4px solid #e94560;
        }
        a { 
            color: #f9b17a; 
            text-decoration: none; 
            font-weight: bold;
        }
        a:hover { 
            color: #e94560; 
            text-decoration: underline;
        }
        h1 {
            text-align: center;
            color: #ff7b25;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎃 Secure Halloween Server - Information</h1>
        
        <div class="endpoint">
            <h3>📋 Available Endpoints:</h3>
            <p><a href="/" target="_blank">/ - Health Check</a> (plain text)</p>
            <p><a href="/halloween" target="_blank">/halloween - Halloween Page</a> (HTML)</p>
            <p><a href="/api/halloween" target="_blank">/api/halloween - JSON API</a> (JSON)</p>
            <p><a href="/info" target="_blank">/info - This Page</a> (HTML)</p>
            <p><a href="/greet?name=Test" target="_blank">/greet?name=... - SECURE Greet</a> (XSS Protected)</p>
        </div>
        
        <div class="endpoint">
            <h3>🛡️ Security Features:</h3>
            <p>✅ Rate Limiting (10 requests/second)</p>
            <p>✅ X-Content-Type-Options: nosniff</p>
            <p>✅ X-Frame-Options: DENY</p>
            <p>✅ X-XSS-Protection: 1; mode=block</p>
            <p>✅ Referrer-Policy: strict-origin-when-cross-origin</p>
            <p>✅ HTML Template Auto-escaping</p>
        </div>
        
        <div class="endpoint">
            <h3>👤 Client Information:</h3>
            <p><strong>IP Address:</strong> ` + r.RemoteAddr + `</p>
            <p><strong>User Agent:</strong> ` + r.UserAgent() + `</p>
            <p><strong>Request Method:</strong> ` + r.Method + `</p>
            <p><strong>Request URL:</strong> ` + r.URL.String() + `</p>
        </div>
    </div>
</body>
</html>`

	_, err := fmt.Fprint(w, info)
	if err != nil {
		log.Printf("Error writing Info response: %v", err)
	}
}
