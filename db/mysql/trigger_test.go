package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMysql_Triggers(t *testing.T) {
	triggers, err := mysql80DB.Triggers("posts")

	require.NoError(t, err)
	require.Len(t, triggers, 1)
	require.Equal(t, "update_posts_updated", triggers[0].Name)
}
