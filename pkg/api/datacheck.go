package api

import (
	"errors"
	"fmt"
	"time"

	database "github.com/PIPILaPUPU/finalTest/pkg/db"
)

func checkDate(task *database.Task) error {
	now := time.Now()

	// Если дата пустая → сегодняшняя
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	// Парсим дату задачи
	date, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return errors.New("wrong date format (must be 20060102)")
	}

	// Если repeat пуст — просто проверяем и при необходимости ставим сегодня
	if task.Repeat == "" {
		if afterNow(now, date) {
			task.Date = now.Format(DateFormat)
		}
		return nil
	}

	// --- Тут важное!!! ---
	// Поддерживаем только: d N   и   y
	// Всё остальное → ошибка

	// Получаем следующую дату через NextDate (он сам проверит формат)
	next, deleteTask, err := NextDate(now, date, task.Repeat)
	if err != nil {
		return errors.New("wrong repeat format")
	}
	if deleteTask {
		return errors.New("repeat rule cannot be empty")
	}

	// Если дата в прошлом — подставляем next
	if afterNow(now, date) {
		task.Date = next.Format(DateFormat)
	}

	return nil
}

func afterNow(now, d time.Time) bool {
	return d.Before(now)
}

func formatID(id int64) string {
	return fmt.Sprintf("%d", id)
}
