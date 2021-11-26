# tigql

![test result](https://github.com/tigql/tigql/actions/workflows/test.yml/badge.svg)

tigql is a graphql engine for tidb and mysql, help expose your database as static or dynamic api easier.

- **static** means you can easily write CRUD code.
- **dynamic** means you can expose database data directly as graphql(like hasura) or rpc call.

## dev
```
# start mysql
docker-compose up mysql57

# install
make install

# start with live reload
make dev

```

## static

```go
// connect to schema testdb
db := mysql.Open("my://root:mypass@localhost:33308/testdb")

// ddl
db.Tables("namePattern") // basic info of tables
db.Table("tableName") // detail info of the table

// data
db.FindAll("posts")
db.FindById("posts", 1)
db.Find("posts", "title = :title and created_at >= :time", name, time)

// native sql
sql := `
    select author, title, created
      from posts
      where author = :author
`
d.Sql(sql, author)
```

## dynamic

visit http://localhost:8080/ for graphql playground, try query

```graphql
query findTables {
  tables {
    name
    type
    comment
    def
    createdAt
  }
}
```

## thanks

- https://github.com/k1LoW/tbls
- https://github.com/directus/directus  
- https://github.com/jeecgboot/jeecg-boot
- https://github.com/Tencent/APIJSON  
