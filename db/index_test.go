package db_test

import (
	"testing"

	"github.com/tigql/tigql/db"

	"github.com/stretchr/testify/require"
)

func TestDB_Indexes(t *testing.T) {
	indexes, err := currentTestDB.Indexes("posts")
	require.NoError(t, err)
	// pk, fk and user defined
	require.Len(t, indexes, 3)
	require.True(t, containsIndex(indexes, "posts_user_id_idx"))
}

func containsIndex(cs []*db.Index, name string) bool {
	for _, c := range cs {
		if c.Name == name {
			return true
		}
	}
	return false
}
