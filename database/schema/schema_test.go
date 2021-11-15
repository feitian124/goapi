package schema

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchema_FindTableByName(t *testing.T) {
	s := Schema{
		Name: "testschema",
		Tables: []*Table{
			{
				Name:    "a",
				Comment: "table a",
			},
			{
				Name:    "b",
				Comment: "table b",
			},
		},
	}
	table, err := s.FindTableByName("b")
	require.NoError(t, err)
	require.Equal(t, "table b", table.Comment)
}

func TestSchema_Sort(t *testing.T) {
	s := Schema{
		Name: "testschema",
		Tables: []*Table{
			{
				Name:    "b",
				Comment: "table b",
			},
			{
				Name:    "a",
				Comment: "table a",
				Columns: []*Column{
					{
						Name:    "b",
						Comment: "column b",
					},
					{
						Name:    "a",
						Comment: "column a",
					},
				},
			},
		},
	}
	err := s.Sort()
	require.NoError(t, err)
	require.Equal(t, "a", s.Tables[0].Name)
	require.Equal(t, "a", s.Tables[0].Columns[0].Name)
}

func TestSchema_Repair(t *testing.T) {
	file, err := os.Open(filepath.Join(testdataDir(), "json_test_schema.json.golden"))
	require.NoError(t, err)

	dec := json.NewDecoder(file)
	got := &Schema{}
	err = dec.Decode(got)
	require.NoError(t, err)
	want := newTestSchema()
	err = got.Repair()
	require.NoError(t, err)

	for i, tt := range got.Tables {
		require.Equal(t, want.Tables[i].Name, got.Tables[i].Name)
		for j := range tt.Columns {
			require.Equal(t, want.Tables[i].Columns[j].Name, got.Tables[i].Columns[j].Name)
			for k := range got.Tables[i].Columns[j].ParentRelations {
				require.Equal(t, want.Tables[i].Columns[j].ParentRelations[k].Table.Name,
					got.Tables[i].Columns[j].ParentRelations[k].Table.Name)
				require.Equal(t, want.Tables[i].Columns[j].ParentRelations[k].ParentTable.Name,
					got.Tables[i].Columns[j].ParentRelations[k].ParentTable.Name)
			}
			for k := range got.Tables[i].Columns[j].ChildRelations {
				require.Equal(t, want.Tables[i].Columns[j].ChildRelations[k].Table.Name,
					got.Tables[i].Columns[j].ChildRelations[k].Table.Name)
				require.Equal(t, want.Tables[i].Columns[j].ChildRelations[k].ParentTable.Name,
					got.Tables[i].Columns[j].ChildRelations[k].ParentTable.Name)
			}
		}
	}

	require.Equal(t, len(want.Relations), len(got.Relations))
}

func testdataDir() string {
	wd, _ := os.Getwd()
	dir, _ := filepath.Abs(filepath.Join(filepath.Dir(wd), "testdata"))
	return dir
}

func newTestSchema() *Schema {
	ca := &Column{
		Name:     "a",
		Type:     "bigint(20)",
		Comment:  "column a",
		Nullable: false,
	}
	cb := &Column{
		Name:     "b",
		Type:     "text",
		Comment:  "column b",
		Nullable: true,
	}

	ta := &Table{
		Name:    "a",
		Type:    "BASE TABLE",
		Comment: "table a",
		Columns: []*Column{
			ca,
			{
				Name:     "a2",
				Type:     "datetime",
				Comment:  "column a2",
				Nullable: false,
				Default: sql.NullString{
					String: "CURRENT_TIMESTAMP",
					Valid:  true,
				},
			},
		},
	}

	tb := &Table{
		Name:    "b",
		Type:    "BASE TABLE",
		Comment: "table b",
		Columns: []*Column{
			cb,
			{
				Name:     "b2",
				Comment:  "column b2",
				Type:     "text",
				Nullable: true,
			},
		},
	}
	r := &Relation{
		Table:         ta,
		Columns:       []*Column{ca},
		ParentTable:   tb,
		ParentColumns: []*Column{cb},
	}
	ca.ParentRelations = []*Relation{r}
	cb.ChildRelations = []*Relation{r}

	s := &Schema{
		Name: "testschema",
		Tables: []*Table{
			ta,
			tb,
		},
		Relations: []*Relation{
			r,
		},
		Driver: &Driver{
			Name:            "testdriver",
			DatabaseVersion: "1.0.0",
		},
	}
	return s
}
