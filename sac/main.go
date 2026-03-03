package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os" // Добавлено для проверки файлов

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// validatePod contains the security logic for pod validation
// It checks if any container is running in privileged mode
// Returns: (allowed bool, message string)
func validatePod(raw []byte) (bool, string) {
	var pod corev1.Pod

	// Parse the pod object from JSON
	if err := json.Unmarshal(raw, &pod); err != nil {
		return false, "Could not decode Pod object"
	}

	// Check each container in the pod
	for _, container := range pod.Spec.Containers {
		// Verify if container has privileged mode enabled
		if container.SecurityContext != nil &&
			container.SecurityContext.Privileged != nil &&
			*container.SecurityContext.Privileged {
			return false, "STOP: Privileged containers are forbidden!"
		}
	}

	return true, "Success: Security check passed"
}

// handleAdmissionReview processes incoming admission webhook requests
func handleAdmissionReview(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Parse the AdmissionReview request
	var ar admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &ar); err != nil {
		http.Error(w, "decoding failed", http.StatusBadRequest)
		return
	}

	// Execute the validation logic on the pod
	allowed, message := validatePod(ar.Request.Object.Raw)

	// Log the validation result
	log.Printf("[WEBHOOK] Request UID: %s | Allowed: %v | Message: %s",
		ar.Request.UID, allowed, message)

	// Build the AdmissionReview response
	responseReview := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     ar.Request.UID,
			Allowed: allowed,
			Result: &metav1.Status{
				Message: message,
			},
		},
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	errResp := json.NewEncoder(w).Encode(responseReview)
	if errResp != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		return
	}
}

func main() {
	// Register the webhook endpoint
	http.HandleFunc("/validate", handleAdmissionReview)

	log.Println("🚀 Webhook server starting on :8443 (HTTPS)")

	// Определяем, какие сертификаты использовать (из секрета или локальные)
	certFile := "cert.pem"
	keyFile := "key.pem"

	// Проверяем, есть ли сертификаты из секрета (приоритет)
	if _, err := os.Stat("/certs/tls.crt"); err == nil {
		certFile = "/certs/tls.crt"
		keyFile = "/certs/tls.key"
		log.Println("📁 Using certificates from /certs/ (Kubernetes secret)")
	} else {
		log.Println("⚠️  Using local certificates from current directory")
	}

	// Start HTTPS server with TLS certificates
	err := http.ListenAndServeTLS(":8443", certFile, keyFile, nil)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
