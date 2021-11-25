package mysql

import (
	"database/sql"

	"github.com/aquasecurity/go-version/pkg/version"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	MinMysqlVersion = "5.7.6"
	DriverName      = "mysql"
	DataSourceName  = "root:mypass@tcp(127.0.0.1:33308)/testdb?parseTime=true"
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
	db                     *sqlx.DB
	supportGeneratedColumn bool
	Schema                 *Schema `json:"schema"`
}

// Open takes a dataSourceName like "root:mypass@tcp(127.0.0.1:33308)/testdb?parseTime=true"
func Open(driverName string, dataSourceName string) (*DB, error) {
	// dsn := "root:mypass@tcp(127.0.0.1:33308)/testdb?parseTime=true"
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	d := &DB{
		Name:   driverName,
		db:     db,
		Schema: &Schema{},
	}
	err = d.CheckVersion()
	if err != nil {
		return nil, err
	}
	err = d.CheckSchema()
	if err != nil {
		return nil, err
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
	row := d.db.QueryRowx(`SELECT version();`)
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

// CheckSchema set schema
func (d *DB) CheckSchema() error {
	var name string
	row := d.db.QueryRowx(`SELECT database();`)
	if err := row.Scan(&name); err != nil {
		return errors.WithStack(err)
	}
	d.Schema.Name = name
	return nil
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// d.logger.Print(query, args...)
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return rows, nil
}

func (d *DB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	// d.logger.Print(query, args...)
	rows, err := d.db.Queryx(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return rows, nil
}

func (d *DB) QueryRowx(query string, args ...interface{}) (*sqlx.Row, error) {
	// d.logger.Print(query, args...)
	rows := d.db.QueryRowx(query, args...)
	return rows, nil
}
