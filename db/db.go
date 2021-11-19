package db

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
)

// DB stands for a connection to database, including current schema
type DB struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Url     string  `json:"url"`
	DB      *sql.DB `json:"db"`
	Schema  *Schema `json:"schema"`
}

func Open(url string) (*DB, error) {
	u, err := dburl.Parse(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	parts := strings.Split(u.Short(), "/")
	if len(parts) < 2 {
		return nil, errors.Errorf("invalid url: parse %s -> %#v", url, u)
	}

	db, err := dburl.Open(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	var d *DB
	if u.Driver == "mysql" {
		d = &DB{
			Name: u.Driver,
			DB:   db,
			Schema: &Schema{
				Name: parts[1],
			},
		}
	}
	return d, nil
}

func (d *DB) Close() error {
	if d != nil && d.DB != nil {
		if err := d.DB.Close(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
