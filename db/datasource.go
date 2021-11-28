package db

import (
	"fmt"

	"github.com/pkg/errors"
)

const CurrentDB = "tidb_5_2"

type Datasource struct {
	UserName string
	Passwd   string
	Host     string
	Port     int
	DBName   string
	Comment  string
}

// ConnectString returns a connection string based on the parameters it's given
// This would normally also contain the password, however we're not using one
func (ds *Datasource) ConnectString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", ds.UserName, ds.Passwd, ds.Host, ds.Port, ds.DBName)
}

func PredefinedDatasource(comment string) (*Datasource, error) {
	tidb52 := &Datasource{UserName: "root", Passwd: "", Host: "192.168.135.154", Port: 4000, DBName: "testdb", Comment: "tidb_5_2"}
	mysql80 := &Datasource{UserName: "root", Passwd: "mypass", Host: "127.0.0.1", Port: 33306, DBName: "testdb", Comment: "mysql_8_0"}
	mysql57 := &Datasource{UserName: "root", Passwd: "mypass", Host: "127.0.0.1", Port: 33307, DBName: "testdb", Comment: "mysql_5_7"}
	mariadb10 := &Datasource{UserName: "root", Passwd: "mypass", Host: "127.0.0.1", Port: 33308, DBName: "testdb", Comment: "mariadb_10_5"}

	testDss := []*Datasource{
		tidb52,
		mysql80,
		mysql57,
		mariadb10,
	}

	for _, ds := range testDss {
		if comment == ds.Comment {
			return ds, nil
		}
	}
	return nil, errors.Errorf("no datasource found by comment %s", comment)
}
