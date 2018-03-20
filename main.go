package main

import (
	"fmt"
	. "gdb/query"
	query "gdb/query/pgsql"
)

func main() {
	builder := query.New()
	builder.From("users")
	builder.Select([]string{"name", "age", "from"})
	builder.Where("name", "young").Quote(func(builder Builder) {
		builder.Where("name", "hackyoung")
		builder.Or()
		builder.Where("name", "hhyoung")
	})
	builder.Where("age", GT, "24")
	builder.WhereIn("name", []string{"h", "w"})
	builder.Quote(func(builder Builder) {
		q := query.New()
		q.From("articles")
		q.Select([]string{"author_id"})
		q.Where("article_name", "时间简史")
		builder.WhereInQuery("id", q)
	})
	builder.OrderBy("name", ASC)
	builder.OrderBy("age", DESC)
	builder.Offset(100).Limit(100)

	fmt.Println(builder.ForQuery())
	fmt.Println(builder.Params())

	builder = query.New()
	builder.From("users")
	builder.Where("name", "yangzhong")
	data := make(map[string]string)
	data["name"] = "yang-zhong"
	data["age"] = "26"

	fmt.Println(builder.ForUpdate(data))
	fmt.Println(builder.Params())

	builder = query.New()
	builder.From("users")
	builder.Where("name", "yangzhong")

	fmt.Println(builder.ForRemove())
	fmt.Println(builder.Params())
}
