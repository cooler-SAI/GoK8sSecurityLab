package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func handleValidate(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Webhook called for validation")

	// Decode the request from Kubernetes
	var admissionReview admissionv1.AdmissionReview
	if err := json.NewDecoder(r.Body).Decode(&admissionReview); err != nil {
		log.Printf("❌ Error decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the pod from the request
	var pod corev1.Pod
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
		log.Printf("❌ Error unmarshaling pod: %v", err)
	}

	// Log information about the pod
	log.Printf("📦 Pod: %s, Namespace: %s", pod.Name, pod.Namespace)

	// Validate the pod (example: reject pods named "bad-pod")
	allowed := true
	message := ""

	if pod.Name == "bad-pod" {
		allowed = false
		message = "Cannot create pod named 'bad-pod'"
		log.Printf("❌ REJECTED: %s", pod.Name)
	} else {
		log.Printf("✅ ALLOWED: %s", pod.Name)
	}

	// Build the response
	admissionResponse := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: allowed,
			Result: &metav1.Status{
				Message: message,
			},
		},
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(admissionResponse); err != nil {
		log.Printf("❌ Error encoding response: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

func main() {
	// Register handlers
	http.HandleFunc("/validate", handleValidate)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "webhooklite is running")
	})

	// Paths to certificates (in the container they will be at /root/certificates/)
	certFile := "certificates/tls.crt"
	keyFile := "certificates/tls.key"

	log.Printf("🔐 HTTPS server starting on port 8443")
	log.Printf("📜 Cert: %s, Key: %s", certFile, keyFile)

	if err := http.ListenAndServeTLS(":8443", certFile, keyFile, nil); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
