package mysql_test

import (
	"log"
	"os"
	"testing"

	"github.com/tigql/tigql/db"
	"github.com/tigql/tigql/db/mysql"
)

var (
	currentTestDatasource *db.Datasource
	currentTestDB         *mysql.DB
)

func TestMain(m *testing.M) {
	var err error

	// if test multiple packages(test ./...), any flags you provide must be valid in all of them.
	// so here use env instead of flag
	comment := os.Getenv("DATASOURCE")
	if len(comment) == 0 {
		comment = db.CurrentDB
	}
	currentTestDatasource, err = db.PredefinedDatasource(comment)
	if err != nil {
		log.Fatal(err)
	}
	currentTestDB, err = mysql.Open(currentTestDatasource)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
