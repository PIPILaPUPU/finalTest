package api

import (
	"encoding/json"
	"net/http"
	"strings"

	database "github.com/PIPILaPUPU/finalTest/pkg/db"
)

// UpdateTaskRequest структура для запроса обновления задачи
type UpdateTaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// updateTaskHandler обрабатывает PUT запрос для обновления задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON запрос
	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if req.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Проверяем правило повторения
	if err := validateRepeatRule(req.Repeat); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	correctedDate, err := isDateValidForTask(req.Date, req.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Создаем объект задачи для обновления
	task := &database.Task{
		ID:      req.ID,
		Date:    correctedDate,
		Title:   strings.TrimSpace(req.Title),
		Comment: strings.TrimSpace(req.Comment),
		Repeat:  strings.TrimSpace(req.Repeat),
	}

	// Обновляем задачу в базе данных
	err = database.UpdateTask(task)
	if err != nil {
		errorMsg := err.Error()
		// Если задача не найдена
		if errorMsg == "задача не найдена" {
			writeJSON(w, map[string]string{"error": errorMsg}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": "Ошибка при обновлении задачи"}, http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем успешный пустой ответ
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
