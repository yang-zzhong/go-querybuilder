# A Database Query Builder

## Feature

1. Support Mysql, Pgsql and easy to support more database
2. Clear Interface 
3. Just Compile Sql, Can Use With All database/sql

## Sample

simple query example

```go
import (
    . "yang-zzhong/database/query"
    builder "yang-zzhong/database/query/mysql"
)

builder := builder.New()
builder.From("users")
builder.Select([]string{"name", "id", "age"})
builder.Where("name", LIKE, "%Frank%")
builder.Quote(func (builder Builder) {
    builder.WhereIn("id", []string{"1", "2", "3"})
    builder.Or().Where("age", GT, "15")
})

// open db

db.Query(builder.ForQuery(), builder.Params()...)

```
query with another table example
```go
users := builder.New().From("users")
users.Select([]string{"id"})
users.Where("name", LIKE, "%Frank%")

articles := builder.New().From("articles")
articles.WhereIn("author_id", users)

// open db
db.Query(articles.ForQuery(), articles.Params()...)

```

update example
```go
users := builder.New().From("users")
users.WhereIn("name", []string{"Stiff", "Chunch"})

data := make(map[string]string)
data["age"] = "50"
data["name"] = "No Name"

// open db
db.Exec(users.ForUpdate(data), users.Params()...)

```

remove example

```go
users := builder.New().From("users")
user.WhereIn("id", []string{"1", "2", "3"})

// open db
db.Exec(users.ForRemove(), users.Params()...)

```
