package querybuilder

import (
	. "testing"
)

func TestForMysql(t *T) {
	builder := NewBuilder(&MysqlModifier{})
	builder.From("users")
	if "SELECT * FROM `users`" != builder.ForQuery() {
		t.Error("From Error")
	}
	builder.Select("name", "age")
	if "SELECT `name`, `age` FROM `users`" != builder.ForQuery() {
		t.Error("Select Error")
	}
	builder = NewBuilder(&MysqlModifier{})
	builder.From("users")
	builder.Select("name", "age")
	if "SELECT `name`, `age` FROM `users` FOR UPDATE" != builder.ForQueryToUpdate() {
		t.Error("Select Error")
	}
	builder = NewBuilder(&MysqlModifier{})
	builder.From("users")
	builder.Where("name", "yang-zzhong")
	if "SELECT * FROM `users` WHERE `name` = ?" != builder.ForQuery() {
		t.Error("Where Error")
	}
	builder.Where("age", GT, 15)
	form := "SELECT * FROM `users` WHERE `name` = ? AND `age` > ?"
	if form != builder.ForQuery() {
		t.Error("Where And Error")
	}
	builder = NewBuilder(&MysqlModifier{})
	builder.From("users")
	builder.Where("name", "yang-zzhong")
	builder.Quote(func(builder *Builder) {
		builder.Where("age", GT, 50)
		builder.Or().Where("age", LT, 10)
	})
	form = "SELECT * FROM `users` WHERE `name` = ? AND (`age` > ? OR `age` < ?)"
	if form != builder.ForQuery() {
		t.Error("QWhere Error")
	}
	data := []map[string]interface{}{}
	for i := 0; i < 10; i++ {
		row := make(map[string]interface{})
		row["name"] = "young"
		row["age"] = 15
		data = append(data, row)
	}
	builder = NewBuilder(&MysqlModifier{})
	builder.From("users").
		Where("name", nil).
		Where("age", NEQ, nil)
	form = "SELECT * FROM `users` WHERE `name` IS NULL AND `age` IS NOT NULL"
	if builder.ForQuery() != form {
		t.Error("Where nil Error")
	}
}

func TestForPostgres(t *T) {
	builder := NewBuilder(&PgsqlModifier{})
	builder.From("users")
	if "SELECT * FROM \"users\"" != builder.ForQuery() {
		t.Error("From Error")
	}
	builder.Select("name", "age")
	if "SELECT \"name\", \"age\" FROM \"users\"" != builder.ForQuery() {
		t.Error("Select Error")
	}
	builder = NewBuilder(&PgsqlModifier{})
	builder.From("users")
	builder.Select("name", "age")
	if "SELECT \"name\", \"age\" FROM \"users\" FOR UPDATE" != builder.ForQueryToUpdate() {
		t.Error("Select Error")
	}
	builder = NewBuilder(&PgsqlModifier{})
	builder.From("users")
	builder.Where("name", "yang-zzhong")
	if "SELECT * FROM \"users\" WHERE \"name\" = $1" != builder.ForQuery() {
		t.Error("Where Error")
	}
	builder.Where("age", GT, 15)
	form := "SELECT * FROM \"users\" WHERE \"name\" = $1 AND \"age\" > $2"
	if form != builder.ForQuery() {
		t.Error("Where And Error")
	}
	builder = NewBuilder(&PgsqlModifier{})
	builder.From("users")
	builder.Where("name", "yang-zzhong")
	builder.Quote(func(builder *Builder) {
		builder.Where("age", GT, 50)
		builder.Or().Where("age", LT, 10)
	})
	form = "SELECT * FROM \"users\" WHERE \"name\" = $1 AND (\"age\" > $2 OR \"age\" < $3)"
	if form != builder.ForQuery() {
		t.Error("QWhere Error")
	}
	data := []map[string]interface{}{}
	for i := 0; i < 10; i++ {
		row := make(map[string]interface{})
		row["name"] = "young"
		row["age"] = 15
		data = append(data, row)
	}
	builder = NewBuilder(&PgsqlModifier{})
	builder.From("users").
		Where("name", nil).
		Where("age", NEQ, nil)
	form = "SELECT * FROM \"users\" WHERE \"name\" IS NULL AND \"age\" IS NOT NULL"
	if builder.ForQuery() != form {
		t.Error("Where nil Error")
	}
}
