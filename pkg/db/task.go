package db

import (
	"database/sql"
	"fmt"
)


type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}


func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query(
		fmt.Sprintf(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT %d`, limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	// Если задач нет, возвращаем пустой слайс, а не nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, rows.Err()
}


// Получить задачу по id
func GetTask(id string) (*Task, error) {
	var t Task
	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Задача не найдена")
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Обновить задачу по id
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}


func DeleteTask(id string) error {
    res, err := DB.Exec(`DELETE FROM scheduler WHERE id=?`, id)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("incorrect id for deleting task")
    }
    return nil
}

func UpdateDate(next string, id string) error {
    res, err := DB.Exec(`UPDATE scheduler SET date=? WHERE id=?`, next, id)
    if err != nil {
        return err
    }
    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("incorrect id for updating date")
    }
    return nil
}
