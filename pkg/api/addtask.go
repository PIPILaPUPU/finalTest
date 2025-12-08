package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	database "github.com/PIPILaPUPU/finalTest/pkg/db"
)

// Структуры для запроса и ответа
type TaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// Вспомогательная функция для отправки JSON
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// Проверка корректности даты
func isValidDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("date is empty")
	}

	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// Проверка, что дата в будущем или сегодня
func isDateValidForTask(taskDate string, repeatRule string) (string, error) {
	now := time.Now()
	today := now.Format("20060102")

	// Если дата не указана, используем сегодняшнюю
	if taskDate == "" {
		return today, nil
	}

	// Проверяем формат даты
	t, err := isValidDate(taskDate)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	// Приводим now к началу дня для корректного сравнения
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Если дата в прошлом
	if t.Before(now) {
		if repeatRule == "" {
			// Без повторения - используем сегодняшнюю дату
			return today, nil
		} else {
			// С повторением - вычисляем следующую дату
			nextDate, deleteFlag, err := NextDate(now, t, repeatRule)
			if err != nil {
				return "", errors.New("неверный формат правила повторения")
			}
			if deleteFlag {
				return "", errors.New("неверный формат правила повторения")
			}
			return nextDate.Format("20060102"), nil
		}
	}

	return taskDate, nil
}

// Проверка правила повторения
func validateRepeatRule(repeatRule string) error {
	if repeatRule == "" {
		return nil
	}

	// Проверяем базовые правила
	repeatRule = strings.TrimSpace(repeatRule)

	if strings.HasPrefix(repeatRule, "d") {
		parts := strings.Split(repeatRule, " ")
		if len(parts) != 2 {
			return errors.New("неверный формат правила повторения")
		}

		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return errors.New("неверный формат правила повторения")
		}
		return nil
	}

	if repeatRule == "y" {
		return nil
	}

	// Для еженедельных и ежемесячных задач - возвращаем ошибку
	if strings.HasPrefix(repeatRule, "w") || strings.HasPrefix(repeatRule, "m") {
		return errors.New("данный тип повторения не поддерживается")
	}

	return errors.New("неверный формат правила повторения")
}

// Основной обработчик добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, TaskResponse{Error: "Ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	// Проверяем заголовок
	if strings.TrimSpace(req.Title) == "" {
		writeJSON(w, TaskResponse{Error: "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Проверяем правило повторения
	if err := validateRepeatRule(req.Repeat); err != nil {
		writeJSON(w, TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	correctedDate, err := isDateValidForTask(req.Date, req.Repeat)
	if err != nil {
		writeJSON(w, TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	// Создаем задачу
	task := &database.Task{
		Date:    correctedDate,
		Title:   strings.TrimSpace(req.Title),
		Comment: strings.TrimSpace(req.Comment),
		Repeat:  strings.TrimSpace(req.Repeat),
	}

	// Добавляем в базу данных
	id, err := database.AddTask(task)
	if err != nil {
		writeJSON(w, TaskResponse{Error: "Ошибка при сохранении задачи"}, http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	writeJSON(w, TaskResponse{ID: strconv.FormatInt(id, 10)}, http.StatusOK)
}
