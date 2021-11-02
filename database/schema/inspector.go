package schema

import "github.com/feitian124/goapi/database/dialects"

type Connection struct {
	dialect  dialects.Dialect
	host     string
	user     string
	password string
	database string
	charset  string
}

type Database struct {
	databaseName string
	Tables       []Table
}
