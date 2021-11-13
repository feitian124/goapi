package mysql

import (
	"fmt"
	"strings"

	"github.com/feitian124/goapi/database/schema"
	"github.com/pkg/errors"
)

const indexSql = `
	SELECT
	(CASE WHEN s.index_name='PRIMARY' AND s.non_unique=0 THEN 'PRIMARY KEY'
		  WHEN s.index_name!='PRIMARY' AND s.non_unique=0 THEN 'UNIQUE KEY'
		  WHEN s.non_unique=1 THEN 'KEY'
		  ELSE null
	  END) AS key_type,
	s.index_name, GROUP_CONCAT(s.column_name ORDER BY s.seq_in_index SEPARATOR ', '), s.index_type
	FROM information_schema.statistics AS s
	LEFT JOIN information_schema.columns AS c ON s.table_schema = c.table_schema 
											 AND s.table_name = c.table_name 
											 AND s.column_name = c.column_name
	WHERE s.table_name = c.table_name
	AND s.table_schema = ?
	AND s.table_name = ?
	GROUP BY key_type, s.table_name, s.index_name, s.index_type
`

// Indexes get a table's indexes
func (m *Mysql) Indexes(schemaName string, tableName string) ([]*schema.Index, error) {
	// indexes
	indexRows, err := m.db.Query(indexSql, schemaName, tableName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer indexRows.Close()

	var indexes []*schema.Index
	for indexRows.Next() {
		var (
			indexKeyType    string
			indexName       string
			indexColumnName string
			indexType       string
			indexDef        string
		)
		err = indexRows.Scan(&indexKeyType, &indexName, &indexColumnName, &indexType)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if indexKeyType == "PRIMARY KEY" {
			indexDef = fmt.Sprintf("%s (%s) USING %s", indexKeyType, indexColumnName, indexType)
		} else {
			indexDef = fmt.Sprintf("%s %s (%s) USING %s", indexKeyType, indexName, indexColumnName, indexType)
		}

		index := &schema.Index{
			Name:    indexName,
			Def:     indexDef,
			Table:   &tableName,
			Columns: strings.Split(indexColumnName, ", "),
		}
		indexes = append(indexes, index)
	}
	return indexes, nil

}
