package api

import ("net/http"
	"encoding/json"
)


// Init регистрирует обработчики API
func Init() {
    http.HandleFunc("/api/nextdate", nextDateHandler)
    http.HandleFunc("/api/task", taskHandler)       // GET, POST, PUT, DELETE
    http.HandleFunc("/api/tasks", tasksHandler)
    http.HandleFunc("/api/task/done", doneTaskHandler) // POST
}




func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}