package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// validatePod — логика безопасности
func validatePod(raw []byte) (bool, string) {
	var pod corev1.Pod
	if err := json.Unmarshal(raw, &pod); err != nil {
		log.Printf("❌ Ошибка декодирования Pod: %v", err)
		return false, "Could not decode Pod object"
	}

	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil &&
			container.SecurityContext.Privileged != nil &&
			*container.SecurityContext.Privileged {
			return false, "STOP: Privileged containers are forbidden!"
		}
	}
	return true, "Success: Security check passed"
}

func handleAdmissionReview(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		log.Printf("❌ Пустой запрос")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var ar admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &ar); err != nil {
		log.Printf("❌ Ошибка парсинга AdmissionReview: %v", err)
		http.Error(w, "decoding failed", http.StatusBadRequest)
		return
	}

	// Важно: проверяем, что Request не nil
	if ar.Request == nil {
		log.Printf("❌ AdmissionReview Request is nil")
		return
	}

	allowed, message := validatePod(ar.Request.Object.Raw)

	log.Printf("[WEBHOOK] UID: %s | Allowed: %v | Message: %s", ar.Request.UID, allowed, message)

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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseReview); err != nil {
		log.Printf("❌ Ошибка кодирования ответа: %v", err)
	}
}

func main() {
	http.HandleFunc("/validate", handleAdmissionReview)

	// Стандартные пути для K8s TLS Secret
	certFile := "/certs/tls.crt"
	keyFile := "/certs/tls.key"

	// Проверка наличия файлов. Если их нет, сервер упадет с логом, а не "тихо"
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Printf("⚠️  Файл %s не найден, пробую локальные...", certFile)
		certFile = "cert.pem"
		keyFile = "key.pem"
	}

	log.Printf("🚀 Webhook запущен на :8443 (HTTPS)")
	log.Printf("📂 Использую сертификат: %s", certFile)

	// Запуск сервера
	server := &http.Server{Addr: ":8443"}
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatalf("❌ КРИТИЧЕСКАЯ ОШИБКА: %v", err)
	}
}
