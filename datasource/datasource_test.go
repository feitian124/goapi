package datasource

import (
	"github.com/stretchr/testify/require"
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
	{config.DSN{URL: "my://root:mypass@localhost:33308/testdb"}, "testdb", 9, 6},
	//{config.DSN{URL: "pg://postgres:pgpass@localhost:55432/testdb?sslmode=disable"}, "testdb", 17, 12},
}

func TestAnalyzeSchema(t *testing.T) {
	for _, tt := range tests {
		schema, err := Analyze(tt.dsn)
		require.NoError(t, err)
		require.Equal(t, schema.Name, tt.schemaName)
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
