package mysql

import (
	"database/sql"
	"strings"
	"time"

	"github.com/tigql/tigql/db"

	"github.com/aquasecurity/go-version/pkg/version"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	MinMysqlVersion = "5.7.6"
	DriverName      = "mysql"
)

// DB stands for a connection to database, including current schema
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to
// close a DB.
type DB struct {
	// Name for database driver, such as mysql
	Name                   string `json:"name"`
	Version                string `json:"version"`
	URL                    string `json:"url"`
	db                     *sqlx.DB
	supportGeneratedColumn bool
	Schema                 *Schema `json:"schema"`
	logger                 *zap.SugaredLogger
}

// Open takes a dataSourceName like "root:mypass@tcp(127.0.0.1:33308)/testdb?parseTime=true"
func Open(ds *db.Datasource) (*DB, error) {
	db, err := sqlx.Open("mysql", ds.ConnectString())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	d := &DB{
		Name:   "mysql",
		db:     db,
		Schema: &Schema{},
	}
	logger, err := InitLogger()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	d.logger = logger
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

func (db *DB) Close() error {
	if db != nil && db.db != nil {
		if err := db.db.Close(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// CheckVersion set version and supportGeneratedColumn
func (db *DB) CheckVersion() error {
	verGeneratedColumn, err := version.Parse(MinMysqlVersion)
	if err != nil {
		return errors.WithStack(err)
	}

	var v string
	row := db.db.QueryRowx(`SELECT version();`)
	if err := row.Scan(&v); err != nil {
		return errors.WithStack(err)
	}
	v = formatMariaVersion(v)
	db.Version = v

	ver, err := version.Parse(v)
	if err != nil {
		return errors.WithStack(err)
	}

	// TODO here meant only for mysql. mariadb and tidb should update
	if ver.LessThan(verGeneratedColumn) {
		db.supportGeneratedColumn = false
	} else {
		db.supportGeneratedColumn = true
	}
	return nil
}

// CheckSchema set schema
func (db *DB) CheckSchema() error {
	var name string
	row := db.db.QueryRowx(`SELECT database();`)
	if err := row.Scan(&name); err != nil {
		return errors.WithStack(err)
	}
	db.Schema.Name = name
	return nil
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now() // 获取当前时间
	db.logger.Debugf("-------------------------------- sql: %s\n", query)
	db.logger.Debugf("---args: %+v\n", args)
	rows, err := db.db.Query(query, args...)
	elapsed := time.Since(start)
	db.logger.Debugf("---time: %s\n", elapsed.String())
	// db.logger.Info("|--rows: %d\n", len(rows))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return rows, nil
}

func (db *DB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	db.logger.Debugf(query, args...)
	rows, err := db.db.Queryx(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return rows, nil
}

func (db *DB) QueryRowx(query string, args ...interface{}) (*sqlx.Row, error) {
	db.logger.Debugf(query, args...)
	rows := db.db.QueryRowx(query, args...)
	return rows, nil
}

func InitLogger() (*zap.SugaredLogger, error) {
	// logger, err := zap.NewDevelopment()
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sugarLogger := logger.Sugar()
	return sugarLogger, nil
}

// formatMariaVersion given 10.5.13-MariaDB-1:10.5.13+maria~focal return 10.5.13
func formatMariaVersion(v string) string {
	if strings.Contains(strings.ToLower(v), "maria") {
		parts := strings.Split(v, "-")
		return parts[0]
	}
	return v
}
