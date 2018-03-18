package mysql

import (
	"regexp"
	. "testing"
	. "yang-zzhong/database/query"
)

func TestForQuery(t *T) {

	builder := New()
	builder.From("users")
	if "SELECT * FROM users" != builder.ForQuery() {
		t.Error("From Error")
	}
	builder.Select([]string{"name", "age"})
	if "SELECT name, age FROM users" != builder.ForQuery() {
		t.Error("Select Error")
	}

	builder = New()
	builder.From("users")
	builder.Where("name", "yang-zzhong")
	if "SELECT * FROM users WHERE name = ?" != builder.ForQuery() {
		t.Error("Where Error")
	}
	builder.Where("age", GT, "15")
	reg := "^SELECT \\* FROM users WHERE name = \\? AND age > \\?$"
	suc, _ := regexp.Match(reg, ([]byte)(builder.ForQuery()))
	if !suc {
		t.Error("Where And Error")
	}

	builder = New()
	builder.From("users")
	builder.Where("name", "yang-zzhong")
	builder.QWhere(func(builder Builder) {
		builder.Where("age", GT, "50")
		builder.Or().Where("age", LT, "10")
	})
	reg = "^SELECT \\* FROM users WHERE name = \\? AND \\(age > \\? OR age < \\?\\)$"
	suc, _ = regexp.Match(reg, ([]byte)(builder.ForQuery()))
	if !suc {
		t.Error("QWhere Error")
	}
}
