package scheduler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, date string, rule string) (string, error) {
	if rule == "" {
		return "", errors.New("пустое правило повторения")
	}

	// Парсим исходную дату
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("некорректный формат даты: %v", err)
	}

	// Разбиваем правило на части
	parts := strings.Fields(rule)
	if len(parts) == 0 {
		return "", errors.New("неверный формат правила")
	}

	command := parts[0]

	switch command {
	case "d":
		// Правило: d <число>
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("некорректное число дней: %v", err)
		}
		if days <= 0 || days > 400 {
			return "", errors.New("число дней должно быть от 1 до 400")
		}
		return handleDailyRule(now, startDate, days)

	case "y":
		// Правило: y
		return handleYearlyRule(now, startDate)

	case "w":
		// Правило: w <дни недели>
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила w")
		}
		return handleWeeklyRule(now, startDate, parts[1])

	case "m":
		// Правило: m <дни месяца> [месяцы]
		if len(parts) < 2 {
			return "", errors.New("неверный формат правила m")
		}
		var monthsStr string
		if len(parts) > 2 {
			monthsStr = parts[2]
		}
		return handleMonthlyRule(now, startDate, parts[1], monthsStr)

	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

// handleDailyRule обрабатывает правило с ежедневным повторением
func handleDailyRule(now, startDate time.Time, days int) (string, error) {
	current := startDate

	// Если начальная дата уже прошла, вычисляем следующую дату
	if !isDateAfter(current, now) {
		// Вычисляем разницу в днях между now и startDate
		daysDiff := int(now.Sub(startDate).Hours() / 24)
		// Вычисляем сколько полных интервалов прошло
		intervalsPassed := daysDiff / days
		// Следующая дата = startDate + (intervalsPassed + 1) * days
		current = startDate.AddDate(0, 0, (intervalsPassed+1)*days)
	}

	return current.Format("20060102"), nil
}

// handleYearlyRule обрабатывает правило с ежегодным повторением
func handleYearlyRule(now, startDate time.Time) (string, error) {
	current := startDate

	// Если начальная дата уже прошла, добавляем годы пока не превысим now
	for !isDateAfter(current, now) {
		// Пытаемся добавить год, сохраняя ту же дату
		nextYear := current.Year() + 1
		nextDate := time.Date(nextYear, current.Month(), current.Day(), 0, 0, 0, 0, current.Location())

		// Если дата невалидна (29 февраля в невисокосном году), берем 1 марта
		if current.Month() == time.February && current.Day() == 29 {
			if !isLeapYear(nextYear) {
				nextDate = time.Date(nextYear, time.March, 1, 0, 0, 0, 0, current.Location())
			}
		}

		current = nextDate
	}

	return current.Format("20060102"), nil
}

// handleWeeklyRule обрабатывает правило с повторением по дням недели
func handleWeeklyRule(now, startDate time.Time, daysStr string) (string, error) {
	// Парсим дни недели
	dayStrs := strings.Split(daysStr, ",")
	weekDays := make([]int, 0, len(dayStrs))

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", fmt.Errorf("некорректный день недели: %v", err)
		}
		if day < 1 || day > 7 {
			return "", errors.New("день недели должен быть от 1 до 7")
		}
		weekDays = append(weekDays, day)
	}

	// Начинаем с начальной даты
	current := startDate

	// Если начальная дата уже прошла, ищем следующую подходящую дату
	if !isDateAfter(current, now) {
		// Начинаем поиск со дня после now
		current = now.AddDate(0, 0, 1)

		// Ищем ближайший подходящий день
		for i := 0; i < 365; i++ {
			currentWeekday := int(current.Weekday())
			if currentWeekday == 0 { // Воскресенье = 7
				currentWeekday = 7
			}

			// Проверяем, подходит ли текущий день
			for _, wd := range weekDays {
				if currentWeekday == wd {
					return current.Format("20060102"), nil
				}
			}
			current = current.AddDate(0, 0, 1)
		}
	} else {
		// Если начальная дата еще не наступила, проверяем подходит ли она
		currentWeekday := int(current.Weekday())
		if currentWeekday == 0 {
			currentWeekday = 7
		}

		for _, wd := range weekDays {
			if currentWeekday == wd {
				return current.Format("20060102"), nil
			}
		}

		// Если не подходит, ищем следующий подходящий день
		for i := 0; i < 365; i++ {
			current = current.AddDate(0, 0, 1)
			currentWeekday := int(current.Weekday())
			if currentWeekday == 0 {
				currentWeekday = 7
			}

			for _, wd := range weekDays {
				if currentWeekday == wd {
					return current.Format("20060102"), nil
				}
			}
		}
	}

	return "", errors.New("не удалось найти следующую дату")
}

// handleMonthlyRule обрабатывает правило с повторением по дням месяца
func handleMonthlyRule(now, startDate time.Time, daysStr, monthsStr string) (string, error) {
	// Парсим дни месяца
	dayStrs := strings.Split(daysStr, ",")
	monthDays := make([]int, 0, len(dayStrs))

	for _, dayStr := range dayStrs {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", fmt.Errorf("некорректный день месяца: %v", err)
		}
		if (day < 1 || day > 31) && day != -1 && day != -2 {
			return "", errors.New("день месяца должен быть от 1 до 31, -1 или -2")
		}
		monthDays = append(monthDays, day)
	}

	// Парсим месяцы (если указаны)
	monthList := make([]int, 0)
	if monthsStr != "" {
		monthStrs := strings.Split(monthsStr, ",")
		for _, monthStr := range monthStrs {
			month, err := strconv.Atoi(strings.TrimSpace(monthStr))
			if err != nil {
				return "", fmt.Errorf("некорректный месяц: %v", err)
			}
			if month < 1 || month > 12 {
				return "", errors.New("месяц должен быть от 1 до 12")
			}
			monthList = append(monthList, month)
		}
	}

	// Начинаем с начальной даты
	current := startDate

	// Если начальная дата уже прошла, ищем следующую подходящую дату
	if !isDateAfter(current, now) {
		// Начинаем поиск со дня после now
		current = now.AddDate(0, 0, 1)
	}

	// Ищем ближайший подходящий день в пределах 2 лет
	for i := 0; i < 730; i++ {
		// Проверяем ограничения по месяцам
		if len(monthList) > 0 {
			monthMatch := false
			currentMonth := int(current.Month())
			for _, m := range monthList {
				if currentMonth == m {
					monthMatch = true
					break
				}
			}
			if !monthMatch {
				current = current.AddDate(0, 0, 1)
				continue
			}
		}

		// Проверяем дни месяца
		for _, md := range monthDays {
			var candidateDay int

			switch md {
			case -1: // Последний день месяца
				candidateDay = lastDayOfMonth(current.Year(), current.Month())
			case -2: // Предпоследний день месяца
				lastDay := lastDayOfMonth(current.Year(), current.Month())
				candidateDay = lastDay - 1
			default:
				candidateDay = md
			}

			// Проверяем валидность дня для данного месяца
			lastDay := lastDayOfMonth(current.Year(), current.Month())
			if candidateDay > lastDay {
				continue
			}

			if current.Day() == candidateDay {
				return current.Format("20060102"), nil
			}
		}
		current = current.AddDate(0, 0, 1)
	}

	return "", errors.New("не удалось найти следующую дату")
}

// lastDayOfMonth возвращает последний день месяца
func lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// isDateAfter проверяет, что date строго после now (игнорируя время)
func isDateAfter(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return dateOnly.After(nowOnly)
}

// isLeapYear проверяет, является ли год високосным
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
