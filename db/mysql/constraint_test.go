package mysql_test

import (
	"testing"

	"github.com/tigql/tigql/db/mysql"

	"github.com/stretchr/testify/require"
)

func TestDB_Constraints(t *testing.T) {
	constraints, err := currentTestDB.Constraints("posts")
	require.NoError(t, err)
	// pk, fk and user defined
	require.Len(t, constraints, 3)
	require.True(t, containsConstraint(constraints, "posts_user_id_fk"))
}

func containsConstraint(cs []*mysql.Constraint, name string) bool {
	for _, c := range cs {
		if c.Name == name {
			return true
		}
	}
	return false
}
