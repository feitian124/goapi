package schema

import (
	"reflect"
	"testing"
)

func TestTable_CollectTablesAndRelations(t1 *testing.T) {
	type fields struct {
		Name             string
		Type             string
		Comment          string
		Columns          []*Column
		Indexes          []*Index
		Constraints      []*Constraint
		Triggers         []*Trigger
		Def              string
		Labels           Labels
		ReferencedTables []*Table
		External         bool
	}
	type args struct {
		distance int
		root     bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Table
		want1   []*Relation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Name:             tt.fields.Name,
				Type:             tt.fields.Type,
				Comment:          tt.fields.Comment,
				Columns:          tt.fields.Columns,
				Indexes:          tt.fields.Indexes,
				Constraints:      tt.fields.Constraints,
				Triggers:         tt.fields.Triggers,
				Def:              tt.fields.Def,
				Labels:           tt.fields.Labels,
				ReferencedTables: tt.fields.ReferencedTables,
				External:         tt.fields.External,
			}
			got, got1, err := t.CollectTablesAndRelations(tt.args.distance, tt.args.root)
			if (err != nil) != tt.wantErr {
				t1.Errorf("CollectTablesAndRelations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("CollectTablesAndRelations() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t1.Errorf("CollectTablesAndRelations() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTable_FindColumnByName(t *testing.T) {
	table := Table{
		Name: "testtable",
		Columns: []*Column{
			{
				Name:    "a",
				Comment: "column a",
			},
			{
				Name:    "b",
				Comment: "column b",
			},
		},
	}
	column, _ := table.FindColumnByName("b")
	want := "column b"
	got := column.Comment
	if got != want {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func TestTable_FindConstrainsByColumnName(t *testing.T) {
	table := Table{
		Name: "testtable",
		Columns: []*Column{
			{
				Name:    "a",
				Comment: "column a",
			},
			{
				Name:    "b",
				Comment: "column b",
			},
		},
	}
	table.Constraints = []*Constraint{
		{
			Name:              "PRIMARY",
			Type:              "PRIMARY KEY",
			Def:               "PRIMARY KEY(a)",
			ReferencedTable:   nil,
			Table:             &table.Name,
			Columns:           []string{"a"},
			ReferencedColumns: []string{},
		},
		{
			Name:              "UNIQUE",
			Type:              "UNIQUE",
			Def:               "UNIQUE KEY a (b)",
			ReferencedTable:   nil,
			Table:             &table.Name,
			Columns:           []string{"b"},
			ReferencedColumns: []string{},
		},
	}

	got := table.FindConstrainsByColumnName("a")
	if want := 1; len(got) != want {
		t.Errorf("got %v\nwant %v", len(got), want)
	}
	if want := "PRIMARY"; got[0].Name != want {
		t.Errorf("got %v\nwant %v", got[0].Name, want)
	}
}

func TestTable_FindConstraintByName(t1 *testing.T) {
	type fields struct {
		Name             string
		Type             string
		Comment          string
		Columns          []*Column
		Indexes          []*Index
		Constraints      []*Constraint
		Triggers         []*Trigger
		Def              string
		Labels           Labels
		ReferencedTables []*Table
		External         bool
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Constraint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Name:             tt.fields.Name,
				Type:             tt.fields.Type,
				Comment:          tt.fields.Comment,
				Columns:          tt.fields.Columns,
				Indexes:          tt.fields.Indexes,
				Constraints:      tt.fields.Constraints,
				Triggers:         tt.fields.Triggers,
				Def:              tt.fields.Def,
				Labels:           tt.fields.Labels,
				ReferencedTables: tt.fields.ReferencedTables,
				External:         tt.fields.External,
			}
			got, err := t.FindConstraintByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t1.Errorf("FindConstraintByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("FindConstraintByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_FindIndexByName(t1 *testing.T) {
	type fields struct {
		Name             string
		Type             string
		Comment          string
		Columns          []*Column
		Indexes          []*Index
		Constraints      []*Constraint
		Triggers         []*Trigger
		Def              string
		Labels           Labels
		ReferencedTables []*Table
		External         bool
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Index
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Name:             tt.fields.Name,
				Type:             tt.fields.Type,
				Comment:          tt.fields.Comment,
				Columns:          tt.fields.Columns,
				Indexes:          tt.fields.Indexes,
				Constraints:      tt.fields.Constraints,
				Triggers:         tt.fields.Triggers,
				Def:              tt.fields.Def,
				Labels:           tt.fields.Labels,
				ReferencedTables: tt.fields.ReferencedTables,
				External:         tt.fields.External,
			}
			got, err := t.FindIndexByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t1.Errorf("FindIndexByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("FindIndexByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_FindTriggerByName(t1 *testing.T) {
	type fields struct {
		Name             string
		Type             string
		Comment          string
		Columns          []*Column
		Indexes          []*Index
		Constraints      []*Constraint
		Triggers         []*Trigger
		Def              string
		Labels           Labels
		ReferencedTables []*Table
		External         bool
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Trigger
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Name:             tt.fields.Name,
				Type:             tt.fields.Type,
				Comment:          tt.fields.Comment,
				Columns:          tt.fields.Columns,
				Indexes:          tt.fields.Indexes,
				Constraints:      tt.fields.Constraints,
				Triggers:         tt.fields.Triggers,
				Def:              tt.fields.Def,
				Labels:           tt.fields.Labels,
				ReferencedTables: tt.fields.ReferencedTables,
				External:         tt.fields.External,
			}
			got, err := t.FindTriggerByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t1.Errorf("FindTriggerByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("FindTriggerByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_HasColumnWithExtraDef(t1 *testing.T) {
	type fields struct {
		Name             string
		Type             string
		Comment          string
		Columns          []*Column
		Indexes          []*Index
		Constraints      []*Constraint
		Triggers         []*Trigger
		Def              string
		Labels           Labels
		ReferencedTables []*Table
		External         bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Table{
				Name:             tt.fields.Name,
				Type:             tt.fields.Type,
				Comment:          tt.fields.Comment,
				Columns:          tt.fields.Columns,
				Indexes:          tt.fields.Indexes,
				Constraints:      tt.fields.Constraints,
				Triggers:         tt.fields.Triggers,
				Def:              tt.fields.Def,
				Labels:           tt.fields.Labels,
				ReferencedTables: tt.fields.ReferencedTables,
				External:         tt.fields.External,
			}
			if got := t.HasColumnWithExtraDef(); got != tt.want {
				t1.Errorf("HasColumnWithExtraDef() = %v, want %v", got, tt.want)
			}
		})
	}
}
