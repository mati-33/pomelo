package data

import (
	"database/sql"
	"fmt"
	"time"
)

type List struct {
	ID      int64
	Name    string
	Created time.Time
	// todo: task count join?
}

func AddList(db *sql.DB, name string) (int64, error) {
	res, err := db.Exec("INSERT INTO lists (name) VALUES (?)", name)
	if err != nil {
		return 0, fmt.Errorf("ListRepo.Add: exec error: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ListsRepo.Add: id error: %v", err)
	}

	return id, nil
}

func GetList(db *sql.DB, id int64) (List, error) {
	var l List

	row := db.QueryRow("SELECT * FROM lists WHERE id = ?", id)
	if err := row.Scan(&l.ID, &l.Name, &l.Created); err != nil {
		if err == sql.ErrNoRows {
			return l, fmt.Errorf("ListsRepo.Get: no list with id %d", id)
		}
		return l, fmt.Errorf("ListsRepo.Get: query error: %v", err)
	}

	return l, nil
}

func GetAllLists(db *sql.DB) ([]List, error) {
	var lists []List

	rows, err := db.Query("SELECT * FROM lists")
	if err != nil {
		return nil, fmt.Errorf("ListsRepo.GetAll: query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l List
		if err := rows.Scan(&l.ID, &l.Name, &l.Created); err != nil {
			return nil, fmt.Errorf("ListsRepo.GetAll: scan error: %v", err)
		}
		lists = append(lists, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListsRepo.GetAll: rows error: %v", err)
	}

	return lists, nil
}

func ModifyList(db *sql.DB, id int64, name string) error {
	res, err := db.Exec("UPDATE lists SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return fmt.Errorf("ListRepo.Modify: exec error for list %d: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ListRepo.Modify: rows affected error for list %d: %v", id, err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("ListRepo.Modify: unexpected number of rows affected for modifying list %d: %v", id, err)
	}

	return nil
}

func DeleteList(db *sql.DB, id int64) error {
	res, err := db.Exec("DELETE FROM lists WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("ListRepo.Delete: exec error for list %d: %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ListsRepo.Delete: rows affected error for list %d: %v", id, err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("ListsRepo.Delete: unexpected number of rows affected for deletying list %d: %v", id, err)
	}

	return nil
}
