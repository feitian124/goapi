package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMysql_Columns(t *testing.T) {
	t.Parallel()
	m, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	columns, err := m.Columns(s.Name, "posts")

	require.NoError(t, err)
	require.Len(t, columns, 7)
	require.Equal(t, columns[0].Name, "id")
}
