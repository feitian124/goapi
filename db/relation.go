package db

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Relation is the struct for table relation
type Relation struct {
	Table         *Table    `json:"table"`
	Columns       []*Column `json:"columns"`
	ParentTable   *Table    `json:"parent_table" yaml:"parentTable"`
	ParentColumns []*Column `json:"parent_columns" yaml:"parentColumns"`
	Def           string    `json:"def"`
	Virtual       bool      `json:"virtual"`
}

const TypeFK = "FOREIGN KEY"

var reFK = regexp.MustCompile(`FOREIGN KEY \((.+)\) REFERENCES ([^\s]+)\s?\((.+)\)`)

// TableRelations get table's relations from its constraints
func (db *DB) TableRelations(tb *Table) ([]*Relation, error) {
	var relations []*Relation
	for _, c := range tb.Constraints {
		if c.Type == TypeFK {
			relation, err := db.Relation(tb, c)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			relations = append(relations, relation)
		}
	}
	return relations, nil
}

// Relation get one relation from one constraint
func (db *DB) Relation(tb *Table, c *Constraint) (*Relation, error) {
	r := &Relation{
		Table: tb,
		Def:   c.Def,
	}
	result := reFK.FindAllStringSubmatch(r.Def, -1)
	if len(result) == 0 || len(result[0]) < 4 {
		return nil, errors.Errorf("can not parse foreign key: %s", r.Def)
	}
	strColumns := strings.Split(result[0][1], ", ")
	strParentTable := result[0][2]
	strParentColumns := strings.Split(result[0][3], ", ")
	for _, c := range strColumns {
		column, err := r.Table.FindColumnByName(c)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		r.Columns = append(r.Columns, column)
		column.ParentRelations = append(column.ParentRelations, r)
	}
	parentTable, err := db.Table(strParentTable)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	r.ParentTable = parentTable
	for _, c := range strParentColumns {
		column, err := parentTable.FindColumnByName(c)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		r.ParentColumns = append(r.ParentColumns, column)
		column.ChildRelations = append(column.ChildRelations, r)
	}
	return r, nil
}
