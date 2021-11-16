# goapi

![test workflow](https://github.com/feitian124/goapi/actions/workflows/test.yml/badge.svg)

goapi expose database as api, both `static` and `dynamic`:

- **static** means you can easily write CRUD code.
- **dynamic** means you can expose database data directly as restful api or rpc call.

## static

TODO

## dynamic

```html
GET /tables
GET /tables/:table
GET /rows/:table
```

## 查询过滤器

- 全匹配查询  
  查询数据没有特殊格式，默认为全匹配查询, 如 name=jack。

  全匹配支持`高级值规则`用法 (查询内容，带有查询规则符号)
  - 小于查询。 查询内容值规则："lt+ 空格 + 内容", 如 `age=lt 60`
  - 小于等于查询。 查询内容值规则："le+ 空格+ 内容"
  - 大于查询。 查询内容值规则："gt+ 空格+ 内容"
  - 大于等于查询。 查询内容值规则："ge+ 空格+ 内容"

- 模糊查询  
  查询数据格式需加星号, 支持前模糊, 后模糊, 前后模糊, 全模糊. 如 name=*jack
  全模糊即查询字段中间带逗号,该查询方式是将查询条件以星号分割, 遍历数组将每个元素作like查询, 再用or拼接

- 包含查询  
  查询数据格式采用逗号分隔, 如 name=jack,jhon， 生成 in 查询

- 不匹配查询  
  查询数据格式需要加叹号前缀, 如 name=!jack

- 范围查询  
  支持数字，时间的范围查询，如 age_begin=18, age_end=60
  {*}_begin： 表示查询范围开始值
  {*}_end:    表示查询范围结束值

## thanks

<https://github.com/directus/directus>  

<https://github.com/jeecgboot/jeecg-boot>  

<https://github.com/Tencent/APIJSON>  
