package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// -----------------блоки "Добавление задачи" и "Получение список задач"----------------
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetTasks(limit int, search string) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	// Проверяем, является ли search датой в формате 02.01.2006
	if search != "" {
		// Пробуем распарсить как дату
		if t, err := time.Parse("02.01.2006", search); err == nil {
			// Преобразуем в формат 20060102
			dateStr := t.Format("20060102")
			query := `SELECT id, date, title, comment, repeat FROM scheduler 
                     WHERE date = ? ORDER BY date, id LIMIT ?`
			rows, err = db.Query(query, dateStr, limit)
		} else {
			// Ищем по подстроке в заголовке или комментарии
			query := `SELECT id, date, title, comment, repeat FROM scheduler 
                     WHERE title LIKE ? OR comment LIKE ? 
                     ORDER BY date, id LIMIT ?`
			searchPattern := "%" + search + "%"
			rows, err = db.Query(query, searchPattern, searchPattern, limit)
		}
	} else {
		// Без поиска - просто все задачи
		query := `SELECT id, date, title, comment, repeat FROM scheduler 
                 ORDER BY date, id LIMIT ?`
		rows, err = db.Query(query, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Возвращаем пустой слайс вместо nil для корректного JSON
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// -----------------блок "Редактирование задач"--------------------------
func UpdateTask(task *Task) error {
	// Преобразуем строковый ID в int64
	taskID, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный идентификатор задачи")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, taskID)
	if err != nil {
		return err
	}

	// Проверяем, была ли обновлена хотя бы одна строка
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func GetTask(id string) (*Task, error) {
	// Преобразуем строковый ID в int64 для базы данных
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("неверный идентификатор задачи")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	row := db.QueryRow(query, taskID)

	var task Task
	var dbID int64
	err = row.Scan(&dbID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, err
	}

	// Преобразуем числовой ID обратно в строку для JSON
	task.ID = strconv.FormatInt(dbID, 10)

	return &task, nil
}

// -----------------блок "Заканчиваем реализацию API"--------------------------
func UpdateDate(id string, newDate string) error {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный идентификатор задачи")
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := db.Exec(query, newDate, taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("неверный идентификатор задачи")
	}

	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := db.Exec(query, taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
