package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func handleValidate(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Webhook called for validation")

	// Читаем всё тело запроса для отладки
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Error reading body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Логируем первые 200 символов тела запроса для отладки
	if len(body) > 0 {
		log.Printf("📦 Request body (first 200 chars): %s", string(body)[:min(200, len(body))])
	} else {
		log.Printf("📦 Request body is EMPTY!")
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	// Декодируем запрос от Kubernetes
	var admissionReview admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		log.Printf("❌ Error decoding JSON: %v", err)
		log.Printf("❌ Raw body: %s", string(body))
		http.Error(w, fmt.Sprintf("JSON decode error: %v", err), http.StatusBadRequest)
		return
	}

	// Проверяем что Request не nil
	if admissionReview.Request == nil {
		log.Printf("❌ AdmissionReview.Request is nil")
		http.Error(w, "AdmissionReview.Request is nil", http.StatusBadRequest)
		return
	}

	// Извлекаем pod из запроса
	var pod corev1.Pod
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
		log.Printf("❌ Error unmarshaling pod: %v", err)
		// Продолжаем, может это не pod?
	}

	// Логируем информацию о запросе
	log.Printf("📦 Request UID: %s", admissionReview.Request.UID)
	log.Printf("📦 Resource: %s/%s", admissionReview.Request.Resource.Resource, admissionReview.Request.Name)
	log.Printf("📦 Operation: %s", admissionReview.Request.Operation)

	if pod.Name != "" {
		log.Printf("📦 Pod: %s, Namespace: %s", pod.Name, pod.Namespace)
	} else {
		log.Printf("📦 Pod name is empty, resource name: %s", admissionReview.Request.Name)
	}

	// Валидация pod (пример: запрещаем pod с именем "bad-pod")
	allowed := true
	message := ""

	// Проверяем имя pod из разных источников
	podName := pod.Name
	if podName == "" && admissionReview.Request.Name != "" {
		podName = admissionReview.Request.Name
	}

	if podName == "bad-pod" {
		allowed = false
		message = "Cannot create pod named 'bad-pod'"
		log.Printf("❌ REJECTED: %s", podName)
	} else {
		log.Printf("✅ ALLOWED: %s", podName)
	}

	// Формируем ответ
	admissionResponse := &admissionv1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: allowed,
	}

	if !allowed {
		admissionResponse.Result = &metav1.Status{
			Message: message,
			Code:    http.StatusForbidden,
		}
	}

	// Создаем AdmissionReview с ответом
	responseReview := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: admissionResponse,
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseReview); err != nil {
		log.Printf("❌ Error encoding response: %v", err)
		http.Error(w, fmt.Sprintf("Response encoding error: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Response sent: allowed=%v", allowed)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "webhooklite is running\n")
	fmt.Fprintf(w, "Endpoints:\n")
	fmt.Fprintf(w, "  /health - health check\n")
	fmt.Fprintf(w, "  /validate - admission webhook\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Регистрируем обработчики
	http.HandleFunc("/validate", handleValidate)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/", handleRoot)

	// Пути к сертификатам (в контейнере они будут в /root/certificates/)
	certFile := "certificates/tls.crt"
	keyFile := "certificates/tls.key"

	log.Printf("🔐 HTTPS server starting on port 8443")
	log.Printf("📜 Cert: %s, Key: %s", certFile, keyFile)
	log.Printf("📡 Endpoints: /health, /validate")

	if err := http.ListenAndServeTLS(":8443", certFile, keyFile, nil); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
