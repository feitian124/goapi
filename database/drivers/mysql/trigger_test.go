package mysql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMysql_Triggers(t *testing.T) {
	m, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	triggers, err := m.Triggers(s.Name, "posts")

	require.NoError(t, err)
	require.Len(t, triggers, 1)
	require.Equal(t, triggers[0].Name, "update_posts_updated")
}
