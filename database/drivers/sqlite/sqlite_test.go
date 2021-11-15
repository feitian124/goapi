package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/feitian124/goapi/database/schema"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xo/dburl"
)

var (
	s  *schema.Schema
	db *sql.DB
)

func TestMain(m *testing.M) {
	s = &schema.Schema{
		Name: "testdb.sqlite3",
	}
	sqliteFilepath := filepath.Join(testdataDir(), "testdb.sqlite3")

	db, _ = dburl.Open(fmt.Sprintf("sq://%s", sqliteFilepath))

	if exit := m.Run(); exit != 0 {
		db.Close()
		os.Exit(exit)
	}

	db.Close()
}

func TestAnalyzeView(t *testing.T) {
	t.Skip("sqlite not support yet")
	driver := New(db)
	err := driver.Analyze(s)
	require.NoError(t, err)
	view, err := s.FindTableByName("post_comments")
	require.NoError(t, err)
	require.NotEmpty(t, view.Def)
}

func TestInfo(t *testing.T) {
	t.Skip("sqlite not support yet")
	driver := New(db)
	d, err := driver.NewDriver()
	if err != nil {
		t.Errorf("%v", err)
	}
	if d.Name != "sqlite" {
		t.Errorf("got %v\nwant %v", d.Name, "sqlite")
	}
	if d.DatabaseVersion == "" {
		t.Errorf("got not empty string.")
	}
}

func TestParseCheckConstraints(t *testing.T) {
	t.Skip("sqlite not support yet")
	table := &schema.Table{
		Name: "check_constraints",
		Columns: []*schema.Column{
			{
				Name: "id",
			},
			{
				Name: "col",
			},
			{
				Name: "brackets",
			},
			{
				Name: "checkcheck",
			},
			{
				Name: "downcase",
			},
			{
				Name: "nl",
			},
		},
	}
	sql := `CREATE TABLE check_constraints (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  col TEXT CHECK(length(col) > 4),
  brackets TEXT UNIQUE NOT NULL CHECK(((length(brackets) > 4))),
  checkcheck TEXT UNIQUE NOT NULL CHECK(length(checkcheck) > 4),
  downcase TEXT UNIQUE NOT NULL check(length(downcase) > 4),
  nl TEXT UNIQUE NOT
    NULL check(length(nl) > 4 OR
      nl != 'ln')
);`
	tableName := "check_constraints"
	want := []*schema.Constraint{
		{
			Name:    "-",
			Type:    "CHECK",
			Def:     "CHECK(length(col) > 4)",
			Table:   &tableName,
			Columns: []string{"col"},
		},
		{
			Name:    "-",
			Type:    "CHECK",
			Def:     "CHECK(((length(brackets) > 4)))",
			Table:   &tableName,
			Columns: []string{"brackets"},
		},
		{
			Name:    "-",
			Type:    "CHECK",
			Def:     "CHECK(length(checkcheck) > 4)",
			Table:   &tableName,
			Columns: []string{"checkcheck"},
		},
		{
			Name:    "-",
			Type:    "CHECK",
			Def:     "check(length(downcase) > 4)",
			Table:   &tableName,
			Columns: []string{"downcase"},
		},
		{
			Name:    "-",
			Type:    "CHECK",
			Def:     "check(length(nl) > 4 OR nl != 'ln')",
			Table:   &tableName,
			Columns: []string{"nl"},
		},
	}
	got := parseCheckConstraints(table, sql)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %#v\nwant: %#v", got, want)
	}
}

func testdataDir() string {
	wd, _ := os.Getwd()
	dir, _ := filepath.Abs(filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata"))
	return dir
}
