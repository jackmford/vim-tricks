package models

import (
	"database/sql"
	"errors"
	"time"
)

type Trick struct {
	Title    string
	Content  string
	LastUsed time.Time
}

type TrickModel struct {
	DB *sql.DB
}

func (m *TrickModel) Insert(title string, content string) (int, error) {
	stmt := `INSERT INTO tricks (title, content, lastused)
	VALUES(?, ?, UTC_TIMESTAMP())`
	result, err := m.DB.Exec(stmt, title, content)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *TrickModel) Get() (*Trick, error) {
	//stmt := `SELECT id, title, content, lastused FROM tricks
	//WHERE lastused < UTC_TIMESTAMP()`
	stmt := `SELECT title, content, lastused FROM tricks WHERE DATE(lastused) = CURDATE();`

	row := m.DB.QueryRow(stmt)
	t := &Trick{}

	err := row.Scan(&t.Title, &t.Content, &t.LastUsed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			stmt = `SELECT title, content, lastused FROM tricks ORDER BY lastused ASC LIMIT 1;`
			row = m.DB.QueryRow(stmt)
			err := row.Scan(&t.Title, &t.Content, &t.LastUsed)
			if err != nil {
				return nil, err
			}
			stmt = `UPDATE tricks SET lastused = UTC_TIMESTAMP() WHERE title = ?`
			_, err = m.DB.Exec(stmt, t.Title)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return t, nil
}
