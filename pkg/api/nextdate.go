package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

// NextDate — твоя функция вычисления следующей даты.
// Здесь просто подключается готовая реализация.
func NextDate(now time.Time, date time.Time, rule string) (time.Time, bool, error) {
	rule = strings.TrimSpace(rule)

	// Пустое правило → удалить
	if rule == "" {
		return time.Time{}, true, nil
	}

	if strings.HasPrefix(rule, "d") {
		parts := strings.Split(rule, " ")
		if len(parts) != 2 {
			return time.Time{}, false, errors.New("wrong d-rule format")
		}

		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return time.Time{}, false, errors.New("d-rule number out of range")
		}

		next := date.AddDate(0, 0, n)

		for next.Before(now) {
			next = next.AddDate(0, 0, n)
		}

		return next, false, nil
	}

	if rule == "y" {
		next := date.AddDate(1, 0, 0)

		for next.Before(now) {
			next = next.AddDate(1, 0, 0)
		}

		return next, false, nil
	}

	return time.Time{}, false, errors.New("unknown repeat rule")
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	rule := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			fmt.Fprint(w, "")
			return
		}
	}

	if dateStr == "" {
		fmt.Fprint(w, "")
		return
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		fmt.Fprint(w, "")
		return
	}

	next, deleteTask, err := NextDate(now, date, rule)
	if err != nil {
		fmt.Fprint(w, "")
		return
	}

	if deleteTask {
		fmt.Fprint(w, "")
		return
	}

	fmt.Fprint(w, next.Format(DateFormat))
}
