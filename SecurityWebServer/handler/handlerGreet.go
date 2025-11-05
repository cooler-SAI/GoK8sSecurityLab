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
	t := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Greeting</title><style>body{text-align: center; font-family: sans-serif; background: #f0f0f0; padding: 50px;}</style></head><body><h1>Hello, %s!</h1><p>This is the insecure handler.</p></body></html>`

	_, err := fmt.Fprintf(w, t, name)
	if err != nil {
		log.Printf("Error writing vulnerable greet response: %v", err)
	}
}

func SecureGreetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Определяем HTML-шаблон (используем %s для name)
	templateHTML := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Greeting</title><style>body{text-align: center; font-family: sans-serif; background: #e8f5e9; padding: 50px; color: #388e3c;}</style></head><body><h1>Hello, %s!</h1><p>This is the SECURE handler.</p></body></html>`

	// Используем html/template для создания шаблона.
	// ParseGlob или ParseFiles лучше для больших проектов,
	// но для одного шаблона используем Must(Parse)
	t, err := template.New("greet").Parse(templateHTML)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Вместо fmt.Fprintf используем t.Execute
	// Go сам экранирует переменную 'name', превращая '<' в '&lt;'
	if err := t.Execute(w, name); err != nil {
		log.Printf("Error executing template: %v", err)
		return
	}
}
