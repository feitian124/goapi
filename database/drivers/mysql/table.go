package mysql

import (
	"fmt"

	"github.com/feitian124/goapi/database/schema"
	"github.com/pkg/errors"
)

const mariaTableSQL = `
	SELECT table_name, table_type, table_comment 
	FROM information_schema.tables 
	WHERE table_schema = ? ORDER BY table_name
`

const mysqlTableSQL = `
	SELECT table_name, table_type, table_comment 
	FROM information_schema.tables 
	WHERE table_schema = ?
`

const mysqlViewSQL = `
	SELECT view_definition FROM information_schema.views
	WHERE table_schema = ?
	AND table_name = ?
`

func (m *Mysql) Table(schemaName string, tableName string, tableType string, tableComment string) (*schema.Table, error) {
	table := &schema.Table{
		Name:    tableName,
		Type:    tableType,
		Comment: tableComment,
	}

	// table definition
	if tableType == "BASE TABLE" {
		tableDefRows, err := m.db.Query(fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		defer tableDefRows.Close()
		for tableDefRows.Next() {
			var (
				tableName string
				tableDef  string
			)
			err := tableDefRows.Scan(&tableName, &tableDef)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			switch {
			case m.showAutoIncrement:
				table.Def = tableDef
			case m.hideAutoIncrement:
				table.Def = reAI.ReplaceAllLiteralString(tableDef, "")
			default:
				table.Def = reAI.ReplaceAllLiteralString(tableDef, " AUTO_INCREMENT=[Redacted by goapi]")
			}
		}
	}

	// view definition
	if tableType == "VIEW" {
		viewDefRows, err := m.db.Query(mysqlViewSQL, schemaName, tableName)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		defer viewDefRows.Close()
		for viewDefRows.Next() {
			var tableDef string
			err := viewDefRows.Scan(&tableDef)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			table.Def = fmt.Sprintf("CREATE VIEW %s AS (%s)", tableName, tableDef)
		}
	}

	return table, nil
}
