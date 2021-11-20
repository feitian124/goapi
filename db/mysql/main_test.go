package mysql_test

import (
	"os"
	"testing"

	"github.com/feitian124/goapi/db/mysql"
)

const mysql80Url = "my://root:mypass@localhost:33308/testdb?parseTime=true"

var mysql80DB *mysql.DB

func TestMain(m *testing.M) {
	var err error
	mysql80DB, err = mysql.Open(mysql80Url)
	if err != nil {
		os.Exit(1)
	}

	if exitCode := m.Run(); exitCode != 0 {
		os.Exit(exitCode)
	}
}
