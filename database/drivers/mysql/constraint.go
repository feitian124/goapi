package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/feitian124/goapi/database/schema"
	"github.com/pkg/errors"
)

const constraintSql = `
	SELECT
	  kcu.constraint_name,
	  sub.costraint_type,
	  GROUP_CONCAT(kcu.column_name ORDER BY kcu.ordinal_position, position_in_unique_constraint SEPARATOR ', ') AS column_name,
	  kcu.referenced_table_name,
	  GROUP_CONCAT(kcu.referenced_column_name ORDER BY kcu.ordinal_position, position_in_unique_constraint SEPARATOR ', ') AS referenced_column_name
	FROM information_schema.key_column_usage AS kcu
	LEFT JOIN information_schema.columns AS c ON kcu.table_schema = c.table_schema AND kcu.table_name = c.table_name AND kcu.column_name = c.column_name
	INNER JOIN
	  (
	   SELECT
	   kcu.table_schema,
	   kcu.table_name,
	   kcu.constraint_name,
	   kcu.column_name,
	   kcu.referenced_table_name,
	   (CASE WHEN kcu.referenced_table_name IS NOT NULL THEN 'FOREIGN KEY'
			WHEN c.column_key = 'PRI' AND kcu.constraint_name = 'PRIMARY' THEN 'PRIMARY KEY'
			WHEN c.column_key = 'PRI' AND kcu.constraint_name != 'PRIMARY' THEN 'UNIQUE'
			WHEN c.column_key = 'UNI' THEN 'UNIQUE'
			WHEN c.column_key = 'MUL' THEN 'UNIQUE'
			ELSE 'UNKNOWN'
	   END) AS costraint_type
	   FROM information_schema.key_column_usage AS kcu
	   LEFT JOIN information_schema.columns AS c ON kcu.table_schema = c.table_schema AND kcu.table_name = c.table_name AND kcu.column_name = c.column_name
	   WHERE kcu.table_name = ?
	   AND kcu.ordinal_position = 1
	  ) AS sub
	ON kcu.constraint_name = sub.constraint_name
	  AND kcu.table_schema = sub.table_schema
	  AND kcu.table_name = sub.table_name
	  AND (kcu.referenced_table_name = sub.referenced_table_name OR (kcu.referenced_table_name IS NULL AND sub.referenced_table_name IS NULL))
	WHERE kcu.table_schema= ?
	  AND kcu.table_name = ?
	GROUP BY kcu.constraint_name, sub.costraint_type, kcu.referenced_table_name
`

func (m *Mysql) Constraints(schemaName string, tableName string) ([]*schema.Constraint, error) {
	constraintRows, err := m.db.Query(constraintSql, tableName, schemaName, tableName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer constraintRows.Close()

	var constraints []*schema.Constraint
	for constraintRows.Next() {
		var (
			constraintName          string
			constraintType          string
			constraintColumnName    string
			constraintRefTableName  sql.NullString
			constraintRefColumnName sql.NullString
			constraintDef           string
		)
		err = constraintRows.Scan(&constraintName, &constraintType, &constraintColumnName, &constraintRefTableName, &constraintRefColumnName)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		switch constraintType {
		case "PRIMARY KEY":
			constraintDef = fmt.Sprintf("PRIMARY KEY (%s)", constraintColumnName)
		case "UNIQUE":
			constraintDef = fmt.Sprintf("UNIQUE KEY %s (%s)", constraintName, constraintColumnName)
		case "FOREIGN KEY":
			constraintType = schema.TypeFK
			constraintDef = fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s (%s)", constraintColumnName, constraintRefTableName.String, constraintRefColumnName.String)
		case "UNKNOWN":
			constraintDef = fmt.Sprintf("UNKNOWN CONSTRAINT (%s) (%s) (%s)", constraintColumnName, constraintRefTableName.String, constraintRefColumnName.String)
		}

		constraint := &schema.Constraint{
			Name:    constraintName,
			Type:    constraintType,
			Def:     constraintDef,
			Table:   &tableName,
			Columns: strings.Split(constraintColumnName, ", "),
		}
		if constraintRefTableName.String != "" {
			constraint.ReferencedTable = &constraintRefTableName.String
			constraint.ReferencedColumns = strings.Split(constraintRefColumnName.String, ", ")
		}

		constraints = append(constraints, constraint)
	}
	return constraints, nil

}
