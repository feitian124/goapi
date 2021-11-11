package mysql

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMysql_Indexes(t *testing.T) {
	m, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	indexes, err := m.Indexes(s.Name, "posts")

	require.NoError(t, err)
	// pk, fk and user defined
	require.Len(t, indexes, 3)
	require.Equal(t, indexes[0].Name, "posts_user_id_idx")
}
