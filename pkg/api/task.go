package api

import (
	"encoding/json"
	"net/http"
	"github.com/Serikkambarov/todo-scheduler/pkg/db"
	"time"
	
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getTaskHandler(w, r)
    case http.MethodPost:
        addTaskHandler(w, r)
    case http.MethodPut:
        updateTaskHandler(w, r)
    case http.MethodDelete:
        deleteTaskHandler(w, r)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}


// --- GET /api/task?id=<id>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, task)
}
// --- Удаление задачи ---
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": "Ошибка удаления задачи"})
		return
	}

	writeJson(w, map[string]string{})
}



// --- PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверяем и корректируем дату 
	
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}


	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{}) // пустой JSON {}
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        writeJson(w, map[string]string{"error": "Не указан идентификатор"})
        return
    }

    task, err := db.GetTask(id)
    if err != nil {
        writeJson(w, map[string]string{"error": "Задача не найдена"})
        return
    }

    if task.Repeat == "" {
        // одноразовая задача, удаляем
        err = db.DeleteTask(id)
        if err != nil {
            writeJson(w, map[string]string{"error": err.Error()})
            return
        }
    } else {
        // периодическая, вычисляем следующую дату
        now := time.Now()
        next, err := NextDate(now, task.Date, task.Repeat)
        if err != nil {
            writeJson(w, map[string]string{"error": err.Error()})
            return
        }
        err = db.UpdateDate(next, id)
        if err != nil {
            writeJson(w, map[string]string{"error": err.Error()})
            return
        }
    }

    writeJson(w, map[string]interface{}{})
}
