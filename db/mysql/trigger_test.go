package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDB_Triggers(t *testing.T) {
	triggers, err := currentTestDB.Triggers("posts")

	require.NoError(t, err)
	require.Len(t, triggers, 1)
	require.Equal(t, "update_posts_updated", triggers[0].Name)
}
