package mysql_test

import (
	"testing"

	"github.com/feitian124/goapi/db/mysql"

	"github.com/stretchr/testify/require"
)

func TestOpenAndClose(t *testing.T) {
	t.Parallel()
	got, err := mysql.Open(currentTestDatasource)
	require.NoError(t, err)
	require.Equal(t, currentTestDatasource.DBName, got.Schema.Name)
	err = got.Close()
	require.NoError(t, err)
}

func TestDB_CheckVersion(t *testing.T) {
	err := currentTestDB.CheckVersion()
	require.NoError(t, err)
}
