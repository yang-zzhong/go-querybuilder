package query

import (
	helpers "github.com/yang-zzhong/go-helpers"
)

const (
	EQ      string = "="
	NEQ     string = "!="
	GT      string = ">"
	GTE     string = ">="
	LT      string = "<"
	LTE     string = "<="
	LIKE    string = "LIKE"
	NULL    string = "IS NULL"
	NOTNULL string = "IS NOT NULL"
	IN      string = "IN"
	NOTIN   string = "NOT IN"
)

type Where interface {
	Id() string
	String() string
	Params() map[string]string
}

type BaseWhere struct {
	Field string
	Op    string
	Value string
	Query Builder
	Array []string

	id     string
	values map[string]string
}

func InitBaseWhere(where *BaseWhere) {
	where.values = make(map[string]string)
	where.id = helpers.RandString(32)
}

func (where *BaseWhere) Params() map[string]string {
	return where.values
}

func (where *BaseWhere) Id() string {
	return where.id
}

func (where *BaseWhere) String() string {
	value := ""
	switch {
	case where.Query != nil:
		value = "(" + where.Query.ForQuery() + ")"
		where.values = where.Query.Params()
	case where.Value != "":
		name := where.name(where.Field)
		value = "@" + name
		where.values[name] = where.Value
	case where.Array != nil:
		value = "("
		length := len(where.Array)
		for i, item := range where.Array {
			name := where.name(where.Field + (string)(i))
			value += "@" + name
			where.values[name] = item
			if i != length-1 {
				value += ", "
			}
		}
		value += ")"
	}
	if where.Op == NULL || where.Op == NOTNULL {
		return where.Field + " " + where.Op
	}
	if where.Op == LIKE {
		return where.Field + " LIKE " + value
	}

	return helpers.Implode([]string{where.Field, where.Op, value}, " ")
}

func (where *BaseWhere) name(field string) string {
	return field + "-" + helpers.RandString(5)
}
