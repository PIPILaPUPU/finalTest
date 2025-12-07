package api

import (
	"encoding/json"
	"net/http"

	database "github.com/PIPILaPUPU/finalTest/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task

	// --- 1. Десериализация JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "json decode error: " + err.Error()})
		return
	}

	// --- 2. Проверка title
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "title is required"})
		return
	}

	// --- 3. Проверка и корректировка даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// --- 4. Запись в БД
	id, err := database.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// --- 5. Ответ
	writeJSON(w, map[string]string{"id": formatID(id)})
}
