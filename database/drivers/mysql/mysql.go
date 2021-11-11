package mysql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/aquasecurity/go-version/pkg/version"
	"github.com/feitian124/goapi/database/ddl"
	"github.com/feitian124/goapi/database/drivers"
	"github.com/feitian124/goapi/database/schema"
	"github.com/pkg/errors"
)

var reFK = regexp.MustCompile(`FOREIGN KEY \((.+)\) REFERENCES ([^\s]+)\s?\((.+)\)`)
var reAI = regexp.MustCompile(` AUTO_INCREMENT=[\d]+`)
var supportGeneratedColumn = true

// Mysql struct
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

// Analyze MySQL database schema
func (m *Mysql) Analyze(s *schema.Schema) error {
	d, err := m.Info()
	if err != nil {
		return errors.WithStack(err)
	}
	s.Driver = d

	if m.mariaMode {
		verGeneratedColumn, err := version.Parse("10.2")
		if err != nil {
			return err
		}
		splitted := strings.Split(s.Driver.DatabaseVersion, "-")
		v, err := version.Parse(splitted[0])
		if err != nil {
			return err
		}
		if v.LessThan(verGeneratedColumn) {
			supportGeneratedColumn = false
		}
	} else {
		verGeneratedColumn, err := version.Parse("5.7.6")
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

	// tables and comments
	tableRows, err := m.db.Query(m.queryForTables(), s.Name)
	if err != nil {
		return errors.WithStack(err)
	}
	defer tableRows.Close()

	var relations []*schema.Relation
	var tables []*schema.Table

	for tableRows.Next() {

		var tableName string
		var tableType string
		var tableComment string
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

	// Relations
	for _, r := range relations {
		result := reFK.FindAllStringSubmatch(r.Def, -1)
		if len(result) == 0 || len(result[0]) < 4 {
			return errors.Errorf("can not parse foreign key: %s", r.Def)
		}
		strColumns := strings.Split(result[0][1], ", ")
		strParentTable := result[0][2]
		strParentColumns := strings.Split(result[0][3], ", ")
		for _, c := range strColumns {
			column, err := r.Table.FindColumnByName(c)
			if err != nil {
				return err
			}
			r.Columns = append(r.Columns, column)
			column.ParentRelations = append(column.ParentRelations, r)
		}
		parentTable, err := s.FindTableByName(strParentTable)
		if err != nil {
			return err
		}
		r.ParentTable = parentTable
		for _, c := range strParentColumns {
			column, err := parentTable.FindColumnByName(c)
			if err != nil {
				return err
			}
			r.ParentColumns = append(r.ParentColumns, column)
			column.ChildRelations = append(column.ChildRelations, r)
		}
	}
	s.Relations = relations

	// referenced tables of view
	for _, t := range s.Tables {
		if t.Type != "VIEW" {
			continue
		}
		for _, rts := range ddl.ParseReferencedTables(t.Def) {
			rt, err := s.FindTableByName(strings.TrimPrefix(rts, fmt.Sprintf("%s.", s.Name)))
			if err != nil {
				rt = &schema.Table{
					Name:     rts,
					External: true,
				}
			}
			t.ReferencedTables = append(t.ReferencedTables, rt)
		}
	}

	return nil
}

// Info return schema.Driver
func (m *Mysql) Info() (*schema.Driver, error) {
	var v string
	row := m.db.QueryRow(`SELECT version();`)
	err := row.Scan(&v)
	if err != nil {
		return nil, err
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
	if str == "NO" {
		return false
	}
	return true
}

func (m *Mysql) queryForTables() string {
	if m.mariaMode {
		return mariaTableSql
	}
	return mysqlTableSql
}
