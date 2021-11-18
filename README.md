# goapi

![test workflow](https://github.com/feitian124/goapi/actions/workflows/test.yml/badge.svg)

goapi expose database as api, both `static` and `dynamic`:

- **static** means you can easily write CRUD code.
- **dynamic** means you can expose database data directly as restful api or rpc call.

## static

```go
m := db.Open("my://root:mypass@localhost:33308/testdb")
s := m.UseSchema("blog")

// ddl
s.Tables("name = :name", name) // basic info of tables
s.Table("tableName") // detail info of the table

// data
s.FindAll("posts")
s.FindById("posts", 1)
s.Find("posts", "title = :title and created_at >= :time", name, time)

// native sql
sql := `
    select author, title, created
      from posts
      where author = :author
`
s.Sql(sql, author)
```

## dynamic

```sh
# ddl
GET /schemas
GET /schemas/:schemaId
GET /schemas/:schemaId/tables
GET /schemas/:schemaId/tables/:tableId

# data
GET /data/:schemaId/:tableId
```

## 查询过滤器(Filtering)

- 全匹配查询  
  ?name=jack 查询数据没有特殊格式为全匹配查询, 如 name=jack。支持高级键规则.

- 模糊查询  
  ?name=*jack 查询数据格式需加星号, 支持前模糊, 后模糊, 前后模糊
  ?name=*ja*ck* 全模糊即查询字段中间带逗号,该查询方式是将查询条件以星号分割, 遍历数组将每个元素作like查询, 再用or拼接

- 包含查询  
  ?name=jack,jhon 查询数据格式采用逗号分隔, 生成 in 查询

- 不匹配查询  
  ?name=!jack 查询数据格式需要加叹号前缀

- 范围查询  
  ?age_begin=18&age_end=60 表示查询 `18 <= age < 60`, 支持数字，时间的范围查询

- 排序, 分页等
  ?limit=10 指定返回记录的数量
  ?offset=10 指定返回记录的开始位置。
  ?page=2&per_page=100 指定第几页，以及每页的记录数。
  ?sortby=name asc 指定返回结果按照哪个属性排序，以及排序顺序。
  ?and=name,age,sex 查询条件之间使用 and 连接, 默认不需要指定, 所有查询条件用 and 连接
  ?or=name,age 查询条件之间使用 or 连接, 再与其他的查询条件进行 and 连接

## 高级键值

查询键中包含查询规则的我们称之为`高级键`(advanced key, AK), 查询值中包含查询规则的我们称之为`高级值`(advanced value, AV).  
如果查询支持从`高级键`抽取查询条件, 我们称该查询支持`高级键规则`. 同理有 `高级值规则` 和 `高级键值规则`.

高级键的组成遵循高级键模式: [原键][键分隔符][键规则符].  
高级值的组成遵循高级值模式: [原值][值分隔符][值规则符].

目前键值分隔符, 键值规则符相同.

考虑到可读性, url友好, 不与常用的json等变量命名冲突, 默认的`键分隔符`和`值分隔符`均为横杠`-`, 可通过 `advanced_key_separator` 和
`advanced_value_separator` 定制.

因为键更加可控, goapi **主要使用高级键规则**, 某些约定俗成的值写法已经通用, 则使用高级值规则, 如上述的模糊, 包含, 不匹配查询.

### 规则符

- eq 等于
- ne 不等于
- lt 小于
- le 小于
- gt 大于
- ge 大于
- flk 前模糊(frontend like)
- blk 后模糊(backend like)
- alk 前后模糊(around like)
- flk 全模糊(full like)

如支持 age 小于查询的高级键为：?age-lt=18, 高级值为: ?age=18-lt


## thanks

<https://github.com/directus/directus>  

<https://github.com/jeecgboot/jeecg-boot>  

<https://github.com/Tencent/APIJSON>  
