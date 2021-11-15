package datasource

import (
	"testing"

	"github.com/feitian124/goapi/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var tests = []struct {
	dsn           config.DSN
	schemaName    string
	tableCount    int
	relationCount int
}{
	{config.DSN{URL: "my://root:mypass@localhost:33308/testdb"}, "testdb", 9, 6},
	// {config.DSN{URL: "pg://postgres:pgpass@localhost:55432/testdb?sslmode=disable"}, "testdb", 17, 12},
}

func TestAnalyzeSchema(t *testing.T) {
	t.Parallel()
	for _, tt := range tests {
		schema, err := Analyze(tt.dsn)
		require.NoError(t, err)
		require.Equal(t, schema.Name, tt.schemaName)
	}
}

func TestAnalyzeTables(t *testing.T) {
	t.Parallel()
	for _, tt := range tests {
		schema, err := Analyze(tt.dsn)
		require.NoError(t, err)
		require.Len(t, schema.Tables, tt.tableCount)
	}
}

func TestAnalyzeRelations(t *testing.T) {
	t.Parallel()
	for _, tt := range tests {
		schema, err := Analyze(tt.dsn)
		require.NoError(t, err)
		require.Len(t, schema.Relations, tt.relationCount)
	}
}
