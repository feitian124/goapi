package mysql_test

import (
	"testing"

	"github.com/feitian124/goapi/db/mysql"
	"github.com/stretchr/testify/require"
)

func TestDB_Tables(t *testing.T) {
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

	d, err := mysql.Open(mysql80Url)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.Tables(tt.args.pattern)
			require.NoError(t, err)
			require.Equal(t, len(got), tt.want)
		})
	}

	err = d.Close()
	require.NoError(t, err)
}

func TestDB_FindTableDDL(t *testing.T) {
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

	d, err := mysql.Open(mysql80Url)

	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.FindTableDDL(tt.args.tableName, tt.args.tableType)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.args.tableName)
			} else {
				require.NoError(t, err)
				require.Contains(t, got, tt.want)
			}
		})
	}

	err = d.Close()
	require.NoError(t, err)
}
