package mysql

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMysql_Columns(t *testing.T) {
	m, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	columns, err := m.Columns(s.Name, "posts")

	require.NoError(t, err)
	require.Len(t, columns, 7)
	require.Equal(t, columns[0].Name, "id")
}
