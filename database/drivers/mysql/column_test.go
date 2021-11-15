package mysql_test

import (
	"testing"

	"github.com/feitian124/goapi/database/drivers/mysql"

	"github.com/stretchr/testify/require"
)

func TestMysql_Columns(t *testing.T) {
	t.Parallel()
	m, err := mysql.New(db)
	if err != nil {
		t.Fatal(err)
	}
	columns, err := m.Columns(s.Name, "posts")

	require.NoError(t, err)
	require.Len(t, columns, 7)
	require.Equal(t, columns[0].Name, "id")
}
