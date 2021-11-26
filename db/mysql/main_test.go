package mysql_test

import (
	"log"
	"os"
	"testing"

	"github.com/feitian124/goapi/db"

	"github.com/feitian124/goapi/db/mysql"
)

var (
	currentTestDatasource *db.Datasource
	currentTestDB         *mysql.DB
)

func TestMain(m *testing.M) {
	var err error

	tidb52 := &db.Datasource{UserName: "root", Passwd: "mypass", Host: "127.0.0.1", Port: 4000, DBName: "testdb", Comment: "tidb 5.2"}
	mysql80 := &db.Datasource{UserName: "root", Passwd: "mypass", Host: "127.0.0.1", Port: 33306, DBName: "testdb", Comment: "mysql 8.0"}
	mysql57 := &db.Datasource{UserName: "root", Passwd: "mypass", Host: "127.0.0.1", Port: 33307, DBName: "testdb", Comment: "mysql 5.7"}
	mariadb10 := &db.Datasource{UserName: "root", Passwd: "mypass", Host: "127.0.0.1", Port: 33308, DBName: "testdb", Comment: "mariadb 10.5"}

	testDss := []*db.Datasource{
		tidb52,
		mysql80,
		mysql57,
		mariadb10,
	}

	for i := 0; i < len(testDss); i++ {
		currentTestDatasource = testDss[i]
		currentTestDB, err = mysql.Open(currentTestDatasource)
		if err != nil {
			log.Fatal(err)
		}
		if exitCode := m.Run(); exitCode != 0 {
			os.Exit(1)
		}
	}
}
