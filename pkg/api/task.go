package api

import (
	"net/http"
	"strconv"

	database "github.com/PIPILaPUPU/finalTest/pkg/db"
)

type TasksResponse struct {
	Tasks []*database.Task `json:"tasks"`
	Error string           `json:"error,omitempty"`
}

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
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
	}
}

// tasksHandler обрабатывает GET запрос для получения списка задач
func returnTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Поддерживаем только GET запросы
	if r.Method != http.MethodGet {
		writeJSON(w, TasksResponse{Error: "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр search из query string
	search := r.URL.Query().Get("search")

	// Ограничение количества задач (по умолчанию 50)
	limit := 50
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Получаем задачи из базы данных
	tasks, err := database.GetTasks(limit, search)
	if err != nil {
		writeJSON(w, TasksResponse{Error: "Ошибка при получении задач"}, http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	response := TasksResponse{Tasks: tasks}
	writeJSON(w, response, http.StatusOK)
}
