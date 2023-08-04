package db

import (
	"database/sql"
	"time"
)

type StorageService interface {
	Save(format string, values ...any) error
	GetByID(query string, id string) *sql.Row

	Many(query string, params ...any) (*sql.Rows, error)
	One(query string, params ...any) *sql.Row
}

type MySql struct {
	*sql.DB
}

func New(url string) (*MySql, error) {
	db, err := sql.Open("mysql", url)

	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return &MySql{
		db,
	}, nil
}

func (m *MySql) Save(format string, values ...any) error {
	_, err := m.Exec(format, values)

	if err != nil {
		return err
	}

	return nil
}

func (m *MySql) GetByID(str string, id string) *sql.Row {
	return m.QueryRow(str, id)
}

func (m *MySql) Many(query string, args ...any) (*sql.Rows, error) {
	return m.Query(query, args)
}

func (m *MySql) One(query string, args ...any) *sql.Row {
	return m.QueryRow(query, args)
}
