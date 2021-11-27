package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/feitian124/goapi/db"

	"github.com/feitian124/goapi/db/mysql"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/feitian124/goapi/graph"
	"github.com/feitian124/goapi/graph/generated"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	var dsn string
	flag.StringVar(&dsn, "dsn", "my://root:mypass@localhost:33308/testdb", "数据库连接符")
	flag.Parse()

	//
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ds, err := db.PredefinedDatasource(db.CurrentDB)
	if err != nil {
		log.Fatal(err)
	}

	db, err := mysql.Open(ds)
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
