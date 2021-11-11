package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMysql_Indexes(t *testing.T) {
	t.Parallel()
	m, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	indexes, err := m.Indexes(s.Name, "posts")

	require.NoError(t, err)
	// pk, fk and user defined
	require.Len(t, indexes, 3)
	require.Equal(t, "posts_user_id_idx", indexes[0].Name)
}
