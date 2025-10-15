package api

import (
	"encoding/json"
	"net/http"
)

// Init регистрирует обработчики API
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler) // GET, POST, PUT, DELETE
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler) // POST
}

// writeJson — универсальная функция для возврата JSON с нужным статусом
func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
	}
}

// writeError — хелпер для отправки ошибки с сообщением и кодом
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJson(w, status, map[string]string{"error": msg})
}

// writeOk — короткий хелпер для успешного ответа {}
func writeOk(w http.ResponseWriter) {
	writeJson(w, http.StatusOK, map[string]any{})
}
