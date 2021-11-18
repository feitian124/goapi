package db

import (
	"database/sql"
)

type DB struct {

}

type Schema struct {
	Name      string      `json:"name"`
	Desc      string      `json:"desc"`
	db *sql.DB
}

type ISchema interface {
	Tables(schema string) ([]*TableInfo, error)
	Table(schema string, table string) (*Table, error)
}

func Open(url string) (*DB, error) {
	return nil, nil
}

func (d *DB) UseSchema(name string) (*Schema, error) {
	return nil, nil
}
