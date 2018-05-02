package querybuilder

import (
	helpers "github.com/yang-zzhong/go-helpers"
	str "strings"
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
	Params() []interface{}
}

type BaseWhere struct {
	Field string
	Op    string
	Value interface{}
	Query *Builder
	Array []string

	id       string
	values   []interface{}
	modifier Modifier
}

func NewW(modifier Modifier) *BaseWhere {
	where := new(BaseWhere)
	where.values = []interface{}{}
	where.id = helpers.RandString(32)
	where.modifier = modifier

	return where
}

func (where *BaseWhere) Params() []interface{} {
	return where.values
}

func (where *BaseWhere) Id() string {
	return where.id
}

func (where *BaseWhere) String() string {
	value := ""
	switch {
	case where.Query != nil:
		sql := where.Query.Replace(false).ForQuery()
		value = "(" + sql + ")"
		where.values = where.Query.Params()
	case where.Value != "":
		value = where.modifier.PrePh()
		where.values = []interface{}{where.Value}
	case where.Array != nil:
		value = "("
		length := len(where.Array)
		for i, item := range where.Array {
			value += where.modifier.PrePh()
			where.values = append(where.values, item)
			if i != length-1 {
				value += ", "
			}
		}
		value += ")"
	}
	if where.Op == NULL || where.Op == NOTNULL {
		return where.modifier.QuoteName(where.Field) + " " + where.Op
	}
	if where.Op == LIKE {
		return where.modifier.QuoteName(where.Field) + " LIKE " + value
	}
	field := where.modifier.QuoteName(where.Field)
	return str.Join([]string{field, where.Op, value}, " ")
}
