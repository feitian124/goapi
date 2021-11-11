package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMysql_Constraints(t *testing.T) {
	m, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	constraints, err := m.Constraints(s.Name, "posts")

	require.NoError(t, err)
	// pk, fk and user defined
	require.Len(t, constraints, 3)
	require.Equal(t, "posts_user_id_fk", constraints[0].Name)
}
