package db

import "database/sql"

type Schema struct {
	Name    string  `json:"name"`
	Comment string  `json:"comment"`
	DB      *sql.DB `json:"db"`
}

type ISchema interface {
	Tables(schema string) ([]*TableInfo, error)
	Table(schema string, table string) (*Table, error)
}

func (s *Schema) Tables(condition string, params ...string) ([]*TableInfo, error) {
	var tis []*TableInfo
	return tis, nil
}

func (s *Schema) Table(tableName string) (*Table, error) {
	t := &Table{}
	return t, nil
}
