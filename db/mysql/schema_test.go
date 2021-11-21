package mysql_test

import (
	"testing"

	"github.com/feitian124/goapi/db/mysql"
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
		{"mysql80 table ddl", args{"post"}, 2, false},
		{"mysql80 table ddl", args{"posts"}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mysql80DB.Tables(tt.args.pattern)
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
		{"mysql80 table posts", "posts", false},
		{"mysql80 table comments", "comments", false},
		{"mysql80 view post_comments", "post_comments", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mysql80DB.Table(tt.want)
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
		tableType mysql.TableType
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"mysql80 table ddl", args{"posts", mysql.BaseTable}, "CREATE TABLE `posts`", false},
		{"mysql80 view ddl", args{"post_comments", mysql.View}, "CREATE VIEW post_comments AS ", false},
		{"mysql80 not exist ddl", args{"a_table_not_exists", mysql.BaseTable}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mysql80DB.FindTableDDL(tt.args.tableName, tt.args.tableType)
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
