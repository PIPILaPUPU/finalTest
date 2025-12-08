package api

import (
	"net/http"

	database "github.com/PIPILaPUPU/finalTest/pkg/db"
)

// getTaskHandler обрабатывает GET запрос для получения задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из query параметров
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	// Получаем задачу из базы данных
	task, err := database.GetTask(id)
	if err != nil {
		errorMsg := err.Error()
		// Если задача не найдена
		if errorMsg == "задача не найдена" {
			writeJSON(w, map[string]string{"error": errorMsg}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": "Ошибка при получении задачи"}, http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем задачу в формате JSON
	writeJSON(w, task, http.StatusOK)
}
