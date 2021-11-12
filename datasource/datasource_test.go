package datasource

import (
	"testing"

	"github.com/feitian124/goapi/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var tests = []struct {
	dsn           config.DSN
	schemaName    string
	tableCount    int
	relationCount int
}{
	{config.DSN{URL: "my://root:mypass@localhost:33306/testdb"}, "testdb", 9, 6},
	{config.DSN{URL: "pg://postgres:pgpass@localhost:55432/testdb?sslmode=disable"}, "testdb", 17, 12},
}

func TestAnalyzeSchema(t *testing.T) {
	for _, tt := range tests {
		schema, err := Analyze(tt.dsn)
		if err != nil {
			t.Errorf("%s", err)
		}
		want := tt.schemaName
		got := schema.Name
		if got != want {
			t.Errorf("got %v\nwant %v", got, want)
		}
	}
}

func TestAnalyzeTables(t *testing.T) {
	for _, tt := range tests {
		schema, err := Analyze(tt.dsn)
		if err != nil {
			t.Errorf("%s", err)
		}
		want := tt.tableCount
		got := len(schema.Tables)
		if got != want {
			t.Errorf("%v: got %v\nwant %v", tt.dsn, got, want)
		}
	}
}

func TestAnalyzeRelations(t *testing.T) {
	for _, tt := range tests {
		schema, err := Analyze(tt.dsn)
		if err != nil {
			t.Errorf("%s", err)
		}
		want := tt.relationCount
		got := len(schema.Relations)
		if got != want {
			t.Errorf("got %v\nwant %v", got, want)
		}
	}
}
