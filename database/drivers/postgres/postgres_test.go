package postgres

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/feitian124/goapi/database/schema"
	_ "github.com/lib/pq"
	"github.com/xo/dburl"
)

var (
	s  *schema.Schema
	db *sql.DB
)

func TestMain(m *testing.M) {
	s = &schema.Schema{
		Name: "testdb",
	}
	db, _ = dburl.Open("pg://postgres:pgpass@localhost:55413/testdb?sslmode=disable")

	if exit := m.Run(); exit != 0 {
		db.Close()
		os.Exit(exit)
	}

	db.Close()
}

func TestAnalyzeView(t *testing.T) {
	t.Skip("postgres not support yet")
	driver := New(db)
	err := driver.Analyze(s)
	require.NoError(t, err)
	view, err := s.FindTableByName("post_comments")
	require.NoError(t, err)
	require.NotEmpty(t, view.Def)
}

func TestExtraDef(t *testing.T) {
	t.Skip("postgres not support yet")
	driver := New(db)
	if err := driver.Analyze(s); err != nil {
		t.Fatal(err)
	}
	tbl, _ := s.FindTableByName("comments")
	{
		c, _ := tbl.FindColumnByName("post_id_desc")
		got := c.ExtraDef
		if want := "GENERATED ALWAYS AS (post_id * '-1'::integer) STORED"; got != want {
			t.Errorf("got %v\nwant %v", got, want)
		}
	}
}

func TestInfo(t *testing.T) {
	t.Skip("postgres not support yet")
	driver := New(db)
	d, err := driver.NewDriver()
	if err != nil {
		t.Errorf("%v", err)
	}
	if d.Name != "postgres" {
		t.Errorf("got %v\nwant %v", d.Name, "postgres")
	}
	if d.DatabaseVersion == "" {
		t.Errorf("got not empty string.")
	}
}
