package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Handler for incoming requests from Kubernetes
func handleValidate(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	if len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	// Decode the AdmissionReview object from Kubernetes
	var admissionReview admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		fmt.Printf("Parsing error: %v\n", err)
		http.Error(w, "Parsing error", http.StatusBadRequest)
		return
	}

	// Extract pod data from the request
	raw := admissionReview.Request.Object.Raw
	pod := corev1.Pod{}
	if err := json.Unmarshal(raw, &pod); err != nil {
		fmt.Printf("Error parsing pod: %v\n", err)
	}

	// SECURITY LOGIC (Charter)
	allowed := true
	message := "The Guard approves this pod."

	// Check: Prohibit running privileged containers
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil && container.SecurityContext.Privileged != nil {
			if *container.SecurityContext.Privileged {
				allowed = false
				message = "SECURITY ERROR: Privileged containers are prohibited by cluster law!"
				break
			}
		}
	}

	// Build response for the API Server
	response := admissionv1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: allowed,
		Result: &metav1.Status{
			Message: message,
		},
	}

	// Send response back
	admissionReview.Response = &response
	res, err := json.Marshal(admissionReview)
	if err != nil {
		fmt.Printf("Error marshalling response: %v\n", err)
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}

func main() {
	// Route for Kubernetes
	http.HandleFunc("/validate", handleValidate)

	port := "8443"
	// Webhook MUST use TLS (HTTPS)
	certFile := "/etc/webhook/certs/tls.crt"
	keyFile := "/etc/webhook/certs/tls.key"

	fmt.Printf("Guard starting duty on port %s...\n", port)

	// Check for certificate existence
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		fmt.Println("WARNING: Certificates not found, starting plain HTTP for testing (not for K8s)")
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			return
		}
	} else {
		err := http.ListenAndServeTLS(":"+port, certFile, keyFile, nil)
		if err != nil {
			fmt.Printf("TLS startup error: %v\n", err)
		}
	}
}
