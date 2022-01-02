package db_test

import (
	"testing"

	"github.com/tigql/tigql/db"

	"github.com/stretchr/testify/require"
)

func TestDB_Tables(t *testing.T) {
	t.Parallel()
	type args struct {
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"post", args{"post"}, 2, false},
		{"posts", args{"posts"}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := currentTestDB.Tables(tt.args.pattern)
			require.NoError(t, err)
			require.Equal(t, len(got), tt.want)
		})
	}
}

func TestDB_Table(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"table posts", "posts", false},
		{"table comments", "comments", false},
		{"view post_comments", "post_comments", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := currentTestDB.Table(tt.want)
			require.NoError(t, err)
			require.Equal(t, got.Name, tt.want)
			require.Greater(t, len(got.Columns), 0)
			if tt.want != "post_comments" {
				require.Greater(t, len(got.Indexes), 0)
				require.Greater(t, len(got.Constraints), 0)
			}
			if tt.want == "posts" {
				require.Greater(t, len(got.Triggers), 0)
			}
			if tt.want == "post_comments" {
				require.Greater(t, len(got.ReferencedTables), 0)
			}
		})
	}
}

func TestDB_FindTableDDL(t *testing.T) {
	t.Parallel()
	type args struct {
		tableName string
		tableType db.TableType
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"table posts", args{"posts", db.BaseTable}, "CREATE TABLE `posts`", false},
		{"view post_comments", args{"post_comments", db.View}, "CREATE VIEW post_comments AS ", false},
		{"table a_table_not_exists", args{"a_table_not_exists", db.BaseTable}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := currentTestDB.FindTableDDL(tt.args.tableName, tt.args.tableType)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.args.tableName)
			} else {
				require.NoError(t, err)
				require.Contains(t, got, tt.want)
			}
		})
	}
}
