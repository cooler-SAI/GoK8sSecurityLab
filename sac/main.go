package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AdmissionReviewHandler — main handler for requests from Kubernetes
func handleAdmissionReview(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Printf("Error closing request body: %v\n", err)
			}
		}(r.Body)
	}

	if len(body) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// Deserialize the request from K8s
	ar := admissionv1.AdmissionReview{}
	if err := json.Unmarshal(body, &ar); err != nil {
		http.Error(w, "could not decode body", http.StatusBadRequest)
		return
	}

	// Object being created (Pod)
	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	if err := json.Unmarshal(raw, &pod); err != nil {
		http.Error(w, "could not decode pod object", http.StatusBadRequest)
		fmt.Printf("Error parsing pod: %v\n", err)
		return
	}

	fmt.Printf("[WEBHOOK] Checking pod: %s in namespace: %s\n", pod.Name, ar.Request.Namespace)

	// SECURITY LOGIC
	allowed := true
	message := "Security check passed"

	for _, container := range pod.Spec.Containers {
		// Block if container wants to be privileged
		if container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
			allowed = false
			message = "STOP: Privileged containers are strictly forbidden in this cluster!"
			break
		}
	}

	// Format response for Kubernetes
	admissionResponse := &admissionv1.AdmissionResponse{
		UID:     ar.Request.UID,
		Allowed: allowed,
		Result: &metav1.Status{
			Message: message,
		},
	}

	response := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: admissionResponse,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "could not encode response", http.StatusInternalServerError)
		fmt.Printf("Error marshalling response: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	write, err := w.Write(resp)
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
		fmt.Printf("Bytes written: %v\n", write)
		return
	}
}

func main() {
	http.HandleFunc("/validate", handleAdmissionReview)

	fmt.Println("Server starting on :8443 (HTTPS required for K8s)")

	// In a real cluster, paths to TLS certificates should be provided here
	// For testing we'll run it like this, but K8s requires TLS certificates
	// signed by the cluster's CA
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		fmt.Println("\nTIP: TLS certificates are required for the webhook to work.")
	}
}
