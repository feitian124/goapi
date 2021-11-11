package mysql

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMysql_Table(t *testing.T) {
	t.Parallel()
	m, err := New(db)
	if err != nil {
		t.Fatal(err)
	}
	table, err := m.Table(s.Name, "posts", "BASE TABLE", "Posts table")
	require.NoError(t, err)
	require.Equal(t, table.Name, "posts")
	require.True(t, strings.Contains(table.Def, "posts"))

	view, err := m.Table(s.Name, "post_comments", "VIEW", "")
	require.NoError(t, err)
	require.Equal(t, view.Name, "post_comments")
	require.True(t, strings.Contains(view.Def, "post_comments"))
}
