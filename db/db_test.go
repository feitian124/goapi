package db_test

import (
	"testing"

	"github.com/feitian124/goapi/db"
	"github.com/stretchr/testify/require"
)

const mysql80Url = "my://root:mypass@localhost:33308/testdb"

func TestOpen(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		url  string
		want *db.DB
	}{
		{"mysql80", mysql80Url, &db.DB{Name: "mysql", Schema: &db.Schema{Name: "testdb"}}},
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

func TestDB_Close(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		url  string
	}{
		{"mysql80", mysql80Url},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := db.Open(tt.url)
			require.NoError(t, err)
			require.Equal(t, d.Schema.Name, "testdb")

			err = d.Close()
			require.NoError(t, err)
		})
	}
}
