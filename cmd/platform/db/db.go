package db

import (
	"database/sql"
	"time"
)

type StorageService interface {
	Save(format string, values ...any) (int64, error)
	GetByID(query string, id int64) (*sql.Row, error)

	Many(query string, params ...any) (*sql.Rows, error)
	One(query string, params ...any) *sql.Row
}

type MySql struct {
	*sql.DB
}

func New(url string) (StorageService, error) {
	db, err := sql.Open("mysql", url)

	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return &MySql{
		db,
	}, nil
}

func (m *MySql) Save(format string, values ...any) (int64, error) {
	result, err := m.Exec(format, values)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *MySql) GetByID(str string, id int64) (*sql.Row, error) {

}

func (m *MySql) Many(query string, args ...any) (*sql.Rows, error) {
	return m.Query(query, args)
}

func (m *MySql) One(query string, args ...any) *sql.Row {
	return m.QueryRow(query, args)
}
