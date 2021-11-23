package mysql_test

import (
	"testing"

	"github.com/feitian124/goapi/db/mysql"

	"github.com/stretchr/testify/require"
)

func TestOpen(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		url  string
		want *mysql.DB
	}{
		{"mysql80", mysql80Url, &mysql.DB{Name: "mysql", Schema: &mysql.Schema{Name: "testdb"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mysql.Open(tt.url)
			require.NoError(t, err)
			require.Equal(t, tt.want.Name, got.Name)
			require.Equal(t, tt.want.Schema.Name, got.Schema.Name)
			// TODO use below equal directly?
			// require.Equal(t, tt.want, got)

			err = got.Close()
			require.NoError(t, err)
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
			d, err := mysql.Open(tt.url)
			require.NoError(t, err)
			require.Equal(t, d.Schema.Name, "testdb")

			err = d.Close()
			require.NoError(t, err)
		})
	}
}

func TestDB_CheckVersion(t *testing.T) {
	err := mysql80DB.CheckVersion()
	require.NoError(t, err)
}
