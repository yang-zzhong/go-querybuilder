package main

import (
	"fmt"
	. "yang-zzhong/database/query"
	query "yang-zzhong/database/query/mysql"
)

func main() {
	builder := query.New()
	builder.From("users")
	builder.Select([]string{"name", "age", "from"})
	builder.Where("name", "young").QWhere(func(builder Builder) {
		builder.Where("name", "hackyoung")
		builder.Or()
		builder.Where("name", "hhyoung")
	})
	builder.Where("age", GT, "24")
	builder.WhereIn("name", []string{"h", "w"})
	builder.QWhere(func(builder Builder) {
		q := query.New()
		q.From("articles")
		q.Select([]string{"author_id"})
		q.Where("article_name", "时间简史")
		builder.WhereInQuery("id", q)
	})
	builder.OrderBy("name", ASC)
	builder.OrderBy("age", DESC)

	fmt.Println(builder.Query())
}
