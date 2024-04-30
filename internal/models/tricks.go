package models

import (
	"database/sql"
	"time"
)

type Trick struct {
	ID       int
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

func (m *TrickModel) Get(id int) (*Trick, error) {
	return nil, nil
}
