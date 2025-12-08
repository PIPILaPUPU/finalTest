package api

import (
	"net/http"

	database "github.com/PIPILaPUPU/finalTest/pkg/db"
)

// deleteTaskHandler обрабатывает DELETE запрос для удаления задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из query параметров
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	// Удаляем задачу из базы данных
	err := database.DeleteTask(id)
	if err != nil {
		errorMsg := err.Error()
		if errorMsg == "задача не найдена" {
			writeJSON(w, map[string]string{"error": errorMsg}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": "Ошибка при удалении задачи"}, http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем успешный пустой ответ
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
