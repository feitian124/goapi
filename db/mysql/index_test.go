package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDB_Indexes(t *testing.T) {
	indexes, err := mysql80DB.Indexes("posts")
	require.NoError(t, err)
	// pk, fk and user defined
	require.Len(t, indexes, 3)
	require.Equal(t, "posts_user_id_idx", indexes[0].Name)
}
