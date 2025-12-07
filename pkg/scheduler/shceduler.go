package repeat

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(currentDate time.Time, rule string) (time.Time, bool, error) {
	rule = strings.TrimSpace(rule)

	// 1. Если правила нет — нужно удалить задачу
	if rule == "" {
		return time.Time{}, true, nil
	}

	// 2. Правило вида "d <число>"
	if strings.HasPrefix(rule, "d ") {
		parts := strings.Split(rule, " ")
		if len(parts) != 2 {
			return time.Time{}, false, errors.New("wrong d-rule format")
		}

		// Парсим число
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return time.Time{}, false, errors.New("cannot parse number in d-rule")
		}
		if n < 1 || n > 400 {
			return time.Time{}, false, errors.New("d-rule number out of range (1..400)")
		}

		// Возвращаем дату + n дней
		return currentDate.AddDate(0, 0, n), false, nil
	}

	// 3. Ежегодное правило "y"
	if rule == "y" {
		return currentDate.AddDate(1, 0, 0), false, nil
	}

	// 4. Если правило неизвестное
	return time.Time{}, false, errors.New("unknown repeat rule")
}
