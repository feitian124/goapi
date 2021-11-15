package schema

import (
	"database/sql"

	"github.com/pkg/errors"
)

// Table is the struct for database table
type Table struct {
	Name             string        `json:"name"`
	Type             string        `json:"type"`
	Comment          string        `json:"comment"`
	Columns          []*Column     `json:"columns"`
	Indexes          []*Index      `json:"indexes"`
	Constraints      []*Constraint `json:"constraints"`
	Triggers         []*Trigger    `json:"triggers"`
	Def              string        `json:"def"`
	Labels           Labels        `json:"labels,omitempty"`
	ReferencedTables []*Table      `json:"referenced_tables,omitempty" yaml:"referencedTables,omitempty"`
	External         bool          `json:"-"` // Table external to the schema
}

// Index is the struct for database index
type Index struct {
	Name    string   `json:"name"`
	Def     string   `json:"def"`
	Table   *string  `json:"table"`
	Columns []string `json:"columns"`
	Comment string   `json:"comment"`
}

// Column is the struct for table column
type Column struct {
	Name            string         `json:"name"`
	Type            string         `json:"type"`
	Nullable        bool           `json:"nullable"`
	Default         sql.NullString `json:"default"`
	Comment         string         `json:"comment"`
	ExtraDef        string         `json:"extra_def,omitempty" yaml:"extraDef,omitempty"`
	ParentRelations []*Relation    `json:"-"`
	ChildRelations  []*Relation    `json:"-"`
}

// Constraint is the struct for database constraint
type Constraint struct {
	Name              string   `json:"name"`
	Type              string   `json:"type"`
	Def               string   `json:"def"`
	Table             *string  `json:"table"`
	ReferencedTable   *string  `json:"referenced_table" yaml:"referencedTable"`
	Columns           []string `json:"columns"`
	ReferencedColumns []string `json:"referenced_columns" yaml:"referencedColumns"`
	Comment           string   `json:"comment"`
}

// Relation is the struct for table relation
type Relation struct {
	Table         *Table    `json:"table"`
	Columns       []*Column `json:"columns"`
	ParentTable   *Table    `json:"parent_table" yaml:"parentTable"`
	ParentColumns []*Column `json:"parent_columns" yaml:"parentColumns"`
	Def           string    `json:"def"`
	Virtual       bool      `json:"virtual"`
}

// Trigger is the struct for database trigger
type Trigger struct {
	Name    string `json:"name"`
	Def     string `json:"def"`
	Comment string `json:"comment"`
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
	cts := []*Constraint{}
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
	tables := []*Table{}
	relations := []*Relation{}
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

	uTables := []*Table{}
	encounteredT := make(map[string]bool)
	for _, t := range tables {
		if !encounteredT[t.Name] {
			encounteredT[t.Name] = true
			uTables = append(uTables, t)
		}
	}

	uRelations := []*Relation{}
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
