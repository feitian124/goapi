package mysql

import (
	"database/sql"
	"regexp"
	"strings"

	"github.com/aquasecurity/go-version/pkg/version"
	"github.com/feitian124/goapi/database/drivers"
	"github.com/feitian124/goapi/database/schema"
	"github.com/pkg/errors"
)

const (
	MinMysqlVersion   = "5.7.6"
	MinMariadbVersion = "10.2"
)

var (
	reAI                   = regexp.MustCompile(` AUTO_INCREMENT=[\d]+`)
	supportGeneratedColumn = true
)

type Mysql struct {
	db        *sql.DB
	mariaMode bool

	// Show AUTO_INCREMENT with increment number
	showAutoIncrement bool

	// Hide the entire AUTO_INCREMENT clause
	hideAutoIncrement bool
}

func ShowAutoIcrrement() drivers.Option {
	return func(d drivers.Driver) error {
		switch d := d.(type) {
		case *Mysql:
			d.showAutoIncrement = true
		}
		return nil
	}
}

func HideAutoIcrrement() drivers.Option {
	return func(d drivers.Driver) error {
		switch d := d.(type) {
		case *Mysql:
			d.hideAutoIncrement = true
		}
		return nil
	}
}

// New return new Mysql
func New(db *sql.DB, opts ...drivers.Option) (*Mysql, error) {
	m := &Mysql{
		db: db,
	}
	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (m *Mysql) checkVersion(s *schema.Schema) error {
	if m.mariaMode {
		verGeneratedColumn, err := version.Parse(MinMariadbVersion)
		if err != nil {
			return errors.WithStack(err)
		}
		parts := strings.Split(s.Driver.DatabaseVersion, "-")
		v, err := version.Parse(parts[0])
		if err != nil {
			return errors.WithStack(err)
		}
		if v.LessThan(verGeneratedColumn) {
			supportGeneratedColumn = false
		}
	} else {
		verGeneratedColumn, err := version.Parse(MinMysqlVersion)
		if err != nil {
			return err
		}
		v, err := version.Parse(s.Driver.DatabaseVersion)
		if err != nil {
			return err
		}
		if v.LessThan(verGeneratedColumn) {
			supportGeneratedColumn = false
		}
	}
	return nil
}

// Analyze MySQL database schema
func (m *Mysql) Analyze(s *schema.Schema) error {
	d, err := m.NewDriver()
	if err != nil {
		return errors.WithStack(err)
	}
	s.Driver = d

	err = m.checkVersion(s)
	if err != nil {
		return errors.WithStack(err)
	}

	// tables and comments
	tableRows, err := m.db.Query(m.queryForTables(), s.Name)
	if err != nil {
		return errors.WithStack(err)
	}
	defer tableRows.Close()

	var tables []*schema.Table

	for tableRows.Next() {

		var tableName, tableType, tableComment string
		err := tableRows.Scan(&tableName, &tableType, &tableComment)
		if err != nil {
			return errors.WithStack(err)
		}

		table, err := m.Table(s.Name, tableName, tableType, tableComment)
		if err != nil {
			return errors.WithStack(err)
		}

		indexes, err := m.Indexes(s.Name, tableName)
		if err != nil {
			return errors.WithStack(err)
		}
		table.Indexes = indexes

		constraints, err := m.Constraints(s.Name, tableName)
		if err != nil {
			return errors.WithStack(err)
		}
		table.Constraints = constraints

		triggers, err := m.Triggers(s.Name, table.Name)
		if err != nil {
			return errors.WithStack(err)
		}
		table.Triggers = triggers

		columns, err := m.Columns(s.Name, table.Name)
		if err != nil {
			return errors.WithStack(err)
		}
		table.Columns = columns

		tables = append(tables, table)
	}

	s.Tables = tables

	err = s.GenRelations()
	if err != nil {
		return errors.WithStack(err)
	}

	s.GenReferencedTables()

	return nil
}

// NewDriver return schema.Driver
func (m *Mysql) NewDriver() (*schema.Driver, error) {
	var v string
	row := m.db.QueryRow(`SELECT version();`)
	err := row.Scan(&v)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	name := "mysql"
	if m.mariaMode {
		name = "mariadb"
	}

	d := &schema.Driver{
		Name:            name,
		DatabaseVersion: v,
	}
	return d, nil
}

// EnableMariaMode enable mariaMode
func (m *Mysql) EnableMariaMode() {
	m.mariaMode = true
}

func convertColumnNullable(str string) bool {
	return str != "NO"
}

func (m *Mysql) queryForTables() string {
	if m.mariaMode {
		return mariaTableSql
	}
	return mysqlTableSql
}
