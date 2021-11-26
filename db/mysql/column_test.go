package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDB_Columns(t *testing.T) {
	t.Parallel()
	columns, err := currentTestDB.Columns("posts")
	require.NoError(t, err)

	require.NoError(t, err)
	require.Len(t, columns, 7)
	require.Equal(t, columns[0].Name, "id")
}
