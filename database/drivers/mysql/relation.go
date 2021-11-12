package mysql

import (
	"strings"

	"github.com/feitian124/goapi/database/schema"
	"github.com/pkg/errors"
)

func getRelationsFromFk(s *schema.Schema, tb *schema.Table) ([]*schema.Relation, error) {
	var relations []*schema.Relation
	for _, c := range tb.Constraints {
		if c.Type == schema.TypeFK {
			relation, err := Relation(s, c)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			relations = append(relations, relation)
		}
	}
	return relations, nil
}

func Relation(s *schema.Schema, c *schema.Constraint) (*schema.Relation, error) {
	r := &schema.Relation{
		Def: c.Def,
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
	parentTable, err := s.FindTableByName(strParentTable)
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
