package data

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      int64
	Name    string
	Created time.Time
	IsDone  bool
	ListID  int64
}

func AddTask(db *sql.DB, listId int64, name string) (int64, error) {
	res, err := db.Exec("INSERT INTO tasks (name, is_done, list_id) VALUES (?, ?, ?)", name, 0, listId)
	if err != nil {
		return 0, fmt.Errorf("AddTask: exec error: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddTask: id error: %v", err)
	}

	return id, nil
}

func GetTask(db *sql.DB, id int64) (Task, error) {
	var t Task

	row := db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
	if err := row.Scan(&t.ID, &t.Name, &t.Created, &t.IsDone, &t.ListID); err != nil {
		if err == sql.ErrNoRows {
			return t, fmt.Errorf("GetTask: no task with id %d", id)
		}
		return t, fmt.Errorf("GetTask: query error: %v", err)
	}

	return t, nil
}

func GetAllListTasks(db *sql.DB, listId int64) ([]Task, error) {
	var tasks []Task

	rows, err := db.Query("SELECT * FROM tasks WHERE list_id = ? ORDER BY created DESC", listId)
	if err != nil {
		return nil, fmt.Errorf("GetAllTasks: query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Created, &t.IsDone, &t.ListID); err != nil {
			return nil, fmt.Errorf("GetAllTasks: scan error: %v", err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllTasks: rows error: %v", err)
	}

	return tasks, nil
}

func ModifyTask(db *sql.DB, id int64, name string, isDone bool) error {
	res, err := db.Exec("UPDATE tasks SET name = ?, is_done = ? WHERE id = ?", name, isDone, id)
	if err != nil {
		return fmt.Errorf("ModifyTask: exec error for task %d: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ModifyTask: rows affected error for task %d: %v", id, err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("ModifyTask: unexpected number of rows affected for modifying task %d: %v", id, err)
	}

	return nil
}

func DeleteTask(db *sql.DB, id int64) error {
	res, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteTask: exec error for task %d: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteTask: rows affected error for task %d: %v", id, err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("DeleteTask: unexpected number of rows affected for deletying task %d: %v", id, err)
	}

	return nil
}
