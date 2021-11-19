package mysql

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Schema is a named collection of tables. A schema can also contain views, indexes, sequences, data types, operators,
// and functions. Schemas are analogous to directories at the operating system level, except schemas cannot be nested.
// For mysql, schema is equal to database. For postgres, a database could have multi schemas.
type Schema struct {
	Name string `json:"name"`
}

const queryTablesByLike = `
	SELECT table_name, table_type, table_comment, create_time
	FROM information_schema.tables
	WHERE table_schema = ?
	AND (table_name like concat('%', ?, '%') or table_comment like concat('%', ?, '%'))
`

const mysqlViewSQL = `
	SELECT view_definition FROM information_schema.views
	WHERE table_schema = ?
	AND table_name = ?
`

type TableType = string

const (
	BaseTable TableType = "BASE TABLE"
	View      TableType = "VIEW"
)

// Tables get table infos using query like "%pattern%" in table name or table comment
func (d *DB) Tables(pattern string) ([]TableInfo, error) {
	var tis []TableInfo
	tableRows, err := d.db.Query(queryTablesByLike, d.Schema.Name, pattern, pattern)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer tableRows.Close()
	for tableRows.Next() {
		var (
			tableName    string
			tableType    string
			tableComment string
			createTime   time.Time
		)
		err := tableRows.Scan(&tableName, &tableType, &tableComment, &createTime)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		ti := TableInfo{
			Name:      tableName,
			Type:      tableType,
			Comment:   tableComment,
			CreatedAt: createTime,
		}

		ti.Def, err = d.FindTableDDL(tableName, tableType)
		if err != nil {
			return nil, err
		}

		tis = append(tis, ti)
	}
	return tis, nil
}

func (d *DB) FindTableDDL(tableName string, tableType TableType) (string, error) {
	if tableType == BaseTable {
		tableDefRows, err := d.db.Query(fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName))
		if err != nil {
			return "", errors.WithStack(err)
		}
		defer tableDefRows.Close()
		for tableDefRows.Next() {
			var (
				tableName string
				tableDef  string
			)
			err := tableDefRows.Scan(&tableName, &tableDef)
			if err != nil {
				return "", errors.WithStack(err)
			}
			return tableDef, nil
		}
	}

	// view definition
	if tableType == View {
		viewDefRows, err := d.db.Query(mysqlViewSQL, d.Schema.Name, tableName)
		if err != nil {
			return "", errors.WithStack(err)
		}
		defer viewDefRows.Close()
		for viewDefRows.Next() {
			var tableDef string
			err := viewDefRows.Scan(&tableDef)
			if err != nil {
				return "", errors.WithStack(err)
			}
			return fmt.Sprintf("CREATE VIEW %s AS (%s)", tableName, tableDef), nil
		}
	}

	return "", nil
}
