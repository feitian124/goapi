package db

import (
	"fmt"

	"github.com/pkg/errors"
)

// Column is the struct for table column
type Column struct {
	Name            string      `json:"name"`
	Type            string      `json:"type"`
	Nullable        bool        `json:"nullable"`
	Default         *string     `json:"default"`
	Comment         *string     `json:"comment"`
	ExtraDef        *string     `json:"extra_def,omitempty" yaml:"extraDef,omitempty"`
	ParentRelations []*Relation `json:"-"`
	ChildRelations  []*Relation `json:"-"`
}

const supportGeneratedColumnSQL = `
	SELECT column_name, column_default, is_nullable, column_type, column_comment, extra, generation_expression
	FROM information_schema.columns
	WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position
`

const columnSQL = `
	SELECT column_name, column_default, is_nullable, column_type, column_comment, extra
	FROM information_schema.columns
	WHERE table_schema = ? AND table_name = ? ORDER BY ordinal_position
`

func (db *DB) Columns(tableName string) ([]*Column, error) {
	columnStmt := supportGeneratedColumnSQL
	if !db.supportGeneratedColumn {
		columnStmt = columnSQL
	}
	columnRows, err := db.Query(columnStmt, db.Schema.Name, tableName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer columnRows.Close()
	var columns []*Column
	for columnRows.Next() {
		var (
			columnName     string
			columnDefault  *string
			isNullable     string
			columnType     string
			columnComment  *string
			extra          *string
			generationExpr *string
		)
		if db.supportGeneratedColumn {
			err = columnRows.Scan(&columnName, &columnDefault, &isNullable, &columnType, &columnComment, &extra, &generationExpr)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			err = columnRows.Scan(&columnName, &columnDefault, &isNullable, &columnType, &columnComment, &extra)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
		extraDef := ""
		if generationExpr != nil {
			switch *extra {
			case "VIRTUAL GENERATED":
				extraDef = fmt.Sprintf("GENERATED ALWAYS AS %s VIRTUAL", *generationExpr)
			case "STORED GENERATED":
				extraDef = fmt.Sprintf("GENERATED ALWAYS AS %s STORED", *generationExpr)
			default:
				extraDef = fmt.Sprintf("%s:%s", extraDef, *generationExpr)
			}
		}
		column := &Column{
			Name:     columnName,
			Type:     columnType,
			Nullable: convertColumnNullable(isNullable),
			Default:  columnDefault,
			Comment:  columnComment,
			ExtraDef: &extraDef,
		}

		columns = append(columns, column)
	}
	return columns, nil
}

func convertColumnNullable(str string) bool {
	return str != "NO"
}
