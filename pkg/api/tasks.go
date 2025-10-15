package api

import (
	"net/http"

	"github.com/Serikkambarov/todo-scheduler/pkg/db"
)

const LIMIT = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(LIMIT)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJson(w, http.StatusOK, TasksResp{Tasks: tasks})
}
