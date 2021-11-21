package mysql

import (
	"time"

	"github.com/pkg/errors"
)

type TableInfo struct {
	Name      string    `json:"name"`
	Type      TableType `json:"type"`
	Comment   string    `json:"comment"`
	Def       string    `json:"def"`
	CreatedAt time.Time `json:"createdAt"`
	External  bool      `json:"-"` // Table external to the schema
}

// Table is the struct for database table
type Table struct {
	TableInfo
	Columns     []*Column     `json:"columns"`
	Indexes     []*Index      `json:"indexes"`
	Constraints []*Constraint `json:"constraints"`
	Triggers    []*Trigger    `json:"triggers"`
	// only used for view
	ReferencedTables []*Table `json:"referenced_tables,omitempty" yaml:"referencedTables,omitempty"`
}

// FindColumnByName find column by column name
func (t *Table) FindColumnByName(name string) (*Column, error) {
	for _, c := range t.Columns {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, errors.Errorf("not found column '%s' on table '%s'", name, t.Name)
}

// FindIndexByName find index by index name
func (t *Table) FindIndexByName(name string) (*Index, error) {
	for _, i := range t.Indexes {
		if i.Name == name {
			return i, nil
		}
	}
	return nil, errors.Errorf("not found index '%s' on table '%s'", name, t.Name)
}

// FindConstraintByName find constraint by constraint name
func (t *Table) FindConstraintByName(name string) (*Constraint, error) {
	for _, c := range t.Constraints {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, errors.Errorf("not found constraint '%s' on table '%s'", name, t.Name)
}

// FindTriggerByName find trigger by trigger name
func (t *Table) FindTriggerByName(name string) (*Trigger, error) {
	for _, trig := range t.Triggers {
		if trig.Name == name {
			return trig, nil
		}
	}
	return nil, errors.Errorf("not found trigger '%s' on table '%s'", name, t.Name)
}

// FindConstrainsByColumnName find constraint by column name
func (t *Table) FindConstrainsByColumnName(name string) []*Constraint {
	var cts []*Constraint
	for _, ct := range t.Constraints {
		for _, ctc := range ct.Columns {
			if ctc == name {
				cts = append(cts, ct)
			}
		}
	}
	return cts
}

func (t *Table) HasColumnWithExtraDef() bool {
	for _, c := range t.Columns {
		if c.ExtraDef != "" {
			return true
		}
	}
	return false
}

func (t *Table) CollectTablesAndRelations(distance int, root bool) ([]*Table, []*Relation, error) {
	var tables []*Table
	var relations []*Relation
	tables = append(tables, t)
	if distance == 0 {
		return tables, relations, nil
	}
	distance--
	for _, c := range t.Columns {
		for _, r := range c.ParentRelations {
			relations = append(relations, r)
			ts, rs, err := r.ParentTable.CollectTablesAndRelations(distance, false)
			if err != nil {
				return nil, nil, err
			}
			tables = append(tables, ts...)
			relations = append(relations, rs...)
		}
		for _, r := range c.ChildRelations {
			relations = append(relations, r)
			ts, rs, err := r.Table.CollectTablesAndRelations(distance, false)
			if err != nil {
				return nil, nil, err
			}
			tables = append(tables, ts...)
			relations = append(relations, rs...)
		}
	}

	if !root {
		return tables, relations, nil
	}

	var uTables []*Table
	encounteredT := make(map[string]bool)
	for _, t := range tables {
		if !encounteredT[t.Name] {
			encounteredT[t.Name] = true
			uTables = append(uTables, t)
		}
	}

	var uRelations []*Relation
	encounteredR := make(map[*Relation]bool)
	for _, r := range relations {
		if !encounteredR[r] {
			encounteredR[r] = true
			if !encounteredT[r.ParentTable.Name] || !encounteredT[r.Table.Name] {
				continue
			}
			uRelations = append(uRelations, r)
		}
	}

	return uTables, uRelations, nil
}
