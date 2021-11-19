package mysql

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
)

// DB stands for a connection to database, including current schema
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to
// close a DB.
type DB struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Url     string  `json:"url"`
	db      *sql.DB `json:"db"`
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

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	var d *DB
	if u.Driver == "mysql" {
		d = &DB{
			Name: u.Driver,
			db:   db,
			Schema: &Schema{
				Name: parts[1],
			},
		}
	}
	return d, nil
}

func (d *DB) Close() error {
	if d != nil && d.db != nil {
		if err := d.db.Close(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
