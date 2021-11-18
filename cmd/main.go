package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/feitian124/goapi/config"
	"github.com/feitian124/goapi/graph"
	"github.com/feitian124/goapi/graph/generated"
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

	//s, err := datasource.Analyze(c.DSN)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(-1)
	//}
	//fmt.Printf("%+v", s)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
