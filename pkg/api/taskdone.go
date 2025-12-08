package api

import (
	"net/http"
	"time"

	database "github.com/PIPILaPUPU/finalTest/pkg/db"
)

// taskDoneHandler обрабатывает POST запрос для отметки выполнения задачи
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID из query параметров
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	// Получаем задачу из базы данных
	task, err := database.GetTask(id)
	if err != nil {
		errorMsg := err.Error()
		if errorMsg == "задача не найдена" {
			writeJSON(w, map[string]string{"error": errorMsg}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": "Ошибка при получении задачи"}, http.StatusInternalServerError)
		}
		return
	}

	// Если правило повторения не указано (пустая строка) - удаляем задачу
	if task.Repeat == "" {
		// Удаляем задачу
		err = database.DeleteTask(id)
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
		return
	}

	// Для периодической задачи вычисляем следующую дату

	// Парсим текущую дату задачи
	currentDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Неверный формат даты в задаче"}, http.StatusInternalServerError)
		return
	}

	// Текущее время для вычисления следующей даты
	now := time.Now()

	// Вычисляем следующую дату
	nextDate, deleteFlag, err := NextDate(now, currentDate, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка при вычислении следующей даты"}, http.StatusInternalServerError)
		return
	}

	// Если NextDate вернула deleteFlag = true, удаляем задачу
	if deleteFlag {
		err = database.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": "Ошибка при удалении задачи"}, http.StatusInternalServerError)
			return
		}
	} else {
		// Обновляем дату задачи
		nextDateStr := nextDate.Format("20060102")
		err = database.UpdateDate(id, nextDateStr)
		if err != nil {
			errorMsg := err.Error()
			if errorMsg == "задача не найдена" {
				writeJSON(w, map[string]string{"error": errorMsg}, http.StatusNotFound)
			} else {
				writeJSON(w, map[string]string{"error": "Ошибка при обновлении даты задачи"}, http.StatusInternalServerError)
			}
			return
		}
	}

	// Возвращаем успешный пустой ответ
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
