package db_test

import (
	"github.com/feitian124/goapi/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOpen(t *testing.T) {
	tests := []struct {
		name    string
		url string
		want    *db.DB
	}{
		{"mysql80", "my://root:mypass@localhost:33308/testdb", &db.DB{ Name:"mysql", Schema: &db.Schema{Name: "testdb"}} },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Open(tt.url)
			require.NoError(t, err)
			require.Equal(t, tt.want.Name, got.Name)
			require.Equal(t, tt.want.Schema.Name, got.Schema.Name)
			// TODO use below equal directly?
			// require.Equal(t, tt.want, got)
		})
	}
}
