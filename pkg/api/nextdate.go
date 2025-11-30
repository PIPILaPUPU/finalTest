package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PIPILaPUPU/finalTest/pkg/scheduler"
)

// NextDateHandler обрабатывает запросы к /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из запроса
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Валидация обязательных параметров
	if dateStr == "" {
		http.Error(w, "Parameter 'date' is required", http.StatusBadRequest)
		return
	}

	if repeat == "" {
		http.Error(w, "Parameter 'repeat' is required", http.StatusBadRequest)
		return
	}

	// Парсим дату now (если не указана - используем текущую)
	var now time.Time
	if nowStr != "" {
		parsedNow, err := time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid 'now' parameter: %v", err), http.StatusBadRequest)
			return
		}
		now = parsedNow
	} else {
		now = time.Now()
	}

	// Вычисляем следующую дату
	nextDate, err := scheduler.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating next date: %v", err), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
