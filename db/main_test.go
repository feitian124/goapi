package db_test

import (
	"log"
	"os"
	"testing"

	"github.com/tigql/tigql/db"
)

var (
	currentTestDatasource *db.Datasource
	currentTestDB         *db.DB
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
	currentTestDB, err = db.Open(currentTestDatasource)
	if err != nil {
		currentTestDB.Close()
		log.Fatal(err)
	}
	os.Exit(m.Run())
	currentTestDB.Close()
}
