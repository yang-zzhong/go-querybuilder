package querybuilder

import (
	"reflect"
)

type column struct {
	fieldName string
	quote     bool
}

type E struct {
	Expr string
}

//
// select columns
//
// builder.Select("name", "age", "id")
// if you wanna same calc on column, you need use struct, for example
// builder.Select(E{"count(1) as number"}, "age").GroupBy("age").ForQuery()
//
func (builder *Builder) Select(cols ...interface{}) *Builder {
	builder.selects = []column{}
	for _, col := range cols {
		value := reflect.ValueOf(col)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
			col = value.Interface()
		}
		switch value.Kind() {
		case reflect.Struct:
			if value.Type().Name() == "E" {
				builder.selects = append(builder.selects, column{col.(E).Expr, false})
			}
		case reflect.String:
			builder.selects = append(builder.selects, column{col.(string), true})
		}
	}
	return builder
}
