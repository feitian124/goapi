package mysql

import (
	"database/sql"
	"strings"

	"github.com/aquasecurity/go-version/pkg/version"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
)

const (
	MinMysqlVersion = "5.7.6"
)

// DB stands for a connection to database, including current schema
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to
// close a DB.
type DB struct {
	Name                   string `json:"name"`
	Version                string `json:"version"`
	URL                    string `json:"url"`
	db                     *sql.DB
	supportGeneratedColumn bool
	Schema                 *Schema `json:"schema"`
}

// Open takes a URL like "protocol+transport://user:pass@host/dbname?option1=a&option2=b"
// for example, a mysql url "my://root:mypass@localhost:33308/testdb?parseTime=true"
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
		err = d.CheckVersion()
		if err != nil {
			return nil, err
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

// CheckVersion set version and supportGeneratedColumn
func (d *DB) CheckVersion() error {
	verGeneratedColumn, err := version.Parse(MinMysqlVersion)
	if err != nil {
		return errors.WithStack(err)
	}

	var v string
	row := d.db.QueryRow(`SELECT version();`)
	if err := row.Scan(&v); err != nil {
		return errors.WithStack(err)
	}
	d.Version = v

	ver, err := version.Parse(v)
	if err != nil {
		return errors.WithStack(err)
	}
	if ver.LessThan(verGeneratedColumn) {
		d.supportGeneratedColumn = false
	} else {
		d.supportGeneratedColumn = true
	}
	return nil
}
