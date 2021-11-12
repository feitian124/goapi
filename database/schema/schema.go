package schema

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/feitian124/goapi/database/ddl"
	"github.com/feitian124/goapi/dict"
	"github.com/pkg/errors"
)

const (
	TypeFK = "FOREIGN KEY"
)

var reFK = regexp.MustCompile(`FOREIGN KEY \((.+)\) REFERENCES ([^\s]+)\s?\((.+)\)`)

type Label struct {
	Name    string
	Virtual bool
}

type Labels []*Label

func (labels Labels) Merge(name string) Labels {
	for _, l := range labels {
		if l.Name == name {
			return labels
		}
	}
	return append(labels, &Label{Name: name, Virtual: true})
}

type DriverMeta struct {
	CurrentSchema string     `json:"current_schema,omitempty" yaml:"currentSchema,omitempty"`
	SearchPaths   []string   `json:"search_paths,omitempty" yaml:"searchPaths,omitempty"`
	Dict          *dict.Dict `json:"dict,omitempty"`
}

// Driver is the struct for database driver information
type Driver struct {
	Name            string      `json:"name"`
	DatabaseVersion string      `json:"database_version" yaml:"databaseVersion"`
	Meta            *DriverMeta `json:"meta"`
}

// Schema is the struct for database schema
type Schema struct {
	Name      string      `json:"name"`
	Desc      string      `json:"desc"`
	Tables    []*Table    `json:"tables"`
	Relations []*Relation `json:"relations"`
	Driver    *Driver     `json:"driver"`
	Labels    Labels      `json:"labels,omitempty"`
}

func (s *Schema) genRelation(tb *Table, c *Constraint) (*Relation, error) {
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

func (s *Schema) genTableRelations(tb *Table) ([]*Relation, error) {
	var relations []*Relation
	for _, c := range tb.Constraints {
		if c.Type == TypeFK {
			relation, err := s.genRelation(tb, c)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			relations = append(relations, relation)
		}
	}
	return relations, nil
}

// GenRelations generate relations from tables
func (s *Schema) GenRelations() error {
	var relations []*Relation
	for _, tb := range s.Tables {
		rs, err := s.genTableRelations(tb)
		if err != nil {
			return errors.WithStack(err)
		}
		relations = append(relations, rs...)
	}
	s.Relations = relations
	return nil
}

// GenReferencedTables generate referenced tables for view
func (s *Schema) GenReferencedTables() {
	for _, tb := range s.Tables {
		if tb.Type != "VIEW" {
			continue
		}
		for _, rts := range ddl.ParseReferencedTables(tb.Def) {
			rt, err := s.FindTableByName(strings.TrimPrefix(rts, fmt.Sprintf("%s.", s.Name)))
			if err != nil {
				rt = &Table{
					Name:     rts,
					External: true,
				}
			}
			tb.ReferencedTables = append(tb.ReferencedTables, rt)
		}
	}
}

func (s *Schema) NormalizeTableName(name string) string {
	if s.Driver != nil && (s.Driver.Name == "postgres" || s.Driver.Name == "redshift") && !strings.Contains(name, ".") {
		return fmt.Sprintf("%s.%s", s.Driver.Meta.CurrentSchema, name)
	}
	return name
}

func (s *Schema) NormalizeTableNames(names []string) []string {
	for i, n := range names {
		names[i] = s.NormalizeTableName(n)
	}
	return names
}

// FindTableByName find table by table name
func (s *Schema) FindTableByName(name string) (*Table, error) {
	for _, t := range s.Tables {
		if s.NormalizeTableName(t.Name) == s.NormalizeTableName(name) {
			return t, nil
		}
	}
	return nil, errors.Errorf("not found table '%s'", name)
}

// FindRelation ...
func (s *Schema) FindRelation(cs, pcs []*Column) (*Relation, error) {
L:
	for _, r := range s.Relations {
		if len(r.Columns) != len(cs) || len(r.ParentColumns) != len(pcs) {
			continue
		}
		for _, rc := range r.Columns {
			exist := false
			for _, cc := range cs {
				if rc == cc {
					exist = true
				}
			}
			if !exist {
				continue L
			}
		}
		for _, rc := range r.ParentColumns {
			exist := false
			for _, cc := range pcs {
				if rc == cc {
					exist = true
				}
			}
			if !exist {
				continue L
			}
		}
		return r, nil
	}
	return nil, errors.Errorf("not found relation '%v, %v'", cs, pcs)
}

// Sort schema tables, columns, relations, and constrains
func (s *Schema) Sort() error {
	for _, t := range s.Tables {
		for _, c := range t.Columns {
			sort.SliceStable(c.ParentRelations, func(i, j int) bool {
				return c.ParentRelations[i].Table.Name < c.ParentRelations[j].Table.Name
			})
			sort.SliceStable(c.ChildRelations, func(i, j int) bool {
				return c.ChildRelations[i].Table.Name < c.ChildRelations[j].Table.Name
			})
		}
		sort.SliceStable(t.Columns, func(i, j int) bool {
			return t.Columns[i].Name < t.Columns[j].Name
		})
		sort.SliceStable(t.Indexes, func(i, j int) bool {
			return t.Indexes[i].Name < t.Indexes[j].Name
		})
		sort.SliceStable(t.Constraints, func(i, j int) bool {
			return t.Constraints[i].Name < t.Constraints[j].Name
		})
		sort.SliceStable(t.Triggers, func(i, j int) bool {
			return t.Triggers[i].Name < t.Triggers[j].Name
		})
	}
	sort.SliceStable(s.Tables, func(i, j int) bool {
		return s.Tables[i].Name < s.Tables[j].Name
	})
	sort.SliceStable(s.Relations, func(i, j int) bool {
		return s.Relations[i].Table.Name < s.Relations[j].Table.Name
	})
	return nil
}

// Repair column relations
func (s *Schema) Repair() error {
	for _, t := range s.Tables {
		if len(t.Columns) == 0 {
			t.Columns = nil
		}
		if len(t.Indexes) == 0 {
			t.Indexes = nil
		}
		if len(t.Constraints) == 0 {
			t.Constraints = nil
		}
		if len(t.Triggers) == 0 {
			t.Triggers = nil
		}
		for i, rt := range t.ReferencedTables {
			tt, err := s.FindTableByName(rt.Name)
			if err != nil {
				rt.External = true
				tt = rt
			}
			t.ReferencedTables[i] = tt
		}
	}

	for _, r := range s.Relations {
		t, err := s.FindTableByName(r.Table.Name)
		if err != nil {
			return errors.Wrap(err, "failed to repair relation")
		}
		for i, rc := range r.Columns {
			c, err := t.FindColumnByName(rc.Name)
			if err != nil {
				return errors.Wrap(err, "failed to repair relation")
			}
			c.ParentRelations = append(c.ParentRelations, r)
			r.Columns[i] = c
		}
		r.Table = t
		pt, err := s.FindTableByName(r.ParentTable.Name)
		if err != nil {
			return errors.Wrap(err, "failed to repair relation")
		}
		for i, rc := range r.ParentColumns {
			pc, err := pt.FindColumnByName(rc.Name)
			if err != nil {
				return errors.Wrap(err, "failed to repair relation")
			}
			pc.ChildRelations = append(pc.ChildRelations, r)
			r.ParentColumns[i] = pc
		}
		r.ParentTable = pt
	}

	return nil
}
