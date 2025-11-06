package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// VulnerableGreetHandler VULNERABLE Handler for XSS demonstration
// WARNING: name is inserted directly into HTML response, making it vulnerable.
func VulnerableGreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// !!! VULNERABLE LINE: name is not sanitized !!!
	simple := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Greeting</title><style>body{text-align: 
center; font-family: sans-serif; background: #f0f0f0; padding: 50px;}</style></head><body><h1>Hello, %s!</h1><p>This is the
insecure handler.</p></body></html>`

	_, err := fmt.Fprintf(w, simple, name)
	if err != nil {
		log.Printf("Error writing vulnerable greet response: %v", err)
	}
}

// SecureGreetHandler - SECURE handler with XSS protection using html/template (EXPORTED)
func SecureGreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	// Secure HTML template with automatic escaping
	const secureTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Secure Greeting</title>
    <style>
        body {
            text-align: center;
            font-family: sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 50px;
            color: white;
        }
        .container {
            background: rgba(255, 255, 255, 0.1);
            padding: 30px;
            border-radius: 15px;
            backdrop-filter: blur(10px);
            max-width: 600px;
            margin: 0 auto;
            border: 1px solid rgba(255, 255, 255, 0.2);
        }
        h1 {
            color: #ffd700;
            margin-bottom: 20px;
        }
        .security-badge {
            background: rgba(76, 175, 80, 0.2);
            padding: 10px;
            border-radius: 8px;
            margin-top: 20px;
            border: 1px solid #4caf50;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Hello, {{.Name}}!</h1>
        <p>This is the <strong>SECURE</strong> handler.</p>
        <div class="security-badge">
            ✅ Protected against XSS attacks<br>
            ✅ Using html/template auto-escaping<br>
            ✅ Safe HTML output
        </div>
        <p><small>Your input has been automatically sanitized for security.</small></p>
    </div>
</body>
</html>`

	// Parse and execute the template with automatic escaping
	tmpl, err := template.New("secureGreet").Parse(secureTemplate)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Execute template with data - html/template automatically escapes HTML characters
	data := struct {
		Name template.HTML
	}{
		Name: template.HTML(name),
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
