package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/feitian124/goapi/config"
	"github.com/feitian124/goapi/datasource"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	var dsn string
	flag.StringVar(&dsn, "dsn", "my://root:mypass@localhost:33308/testdb", "数据库连接符")
	flag.Parse()

	c := config.New()

	if len(dsn) > 0 {
		c.DSN = config.DSN{
			URL: dsn,
		}
	}

	s, err := datasource.Analyze(c.DSN)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("%+v", s)
}
