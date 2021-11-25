package mysql_test

import (
	"log"
	"os"
	"testing"

	"github.com/feitian124/goapi/db/mysql"
)

var mysql80DB *mysql.DB

func TestMain(m *testing.M) {
	var err error
	mysql80DB, err = mysql.Open(mysql.DriverName, mysql.DataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if exitCode := m.Run(); exitCode != 0 {
		os.Exit(exitCode)
	}
}
