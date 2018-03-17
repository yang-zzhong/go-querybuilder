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
	NULL    string = "NULL"
	NOTNULL string = "NOT NULL"
	IN      string = "IN"
	NOTIN   string = "NOT IN"
)

type Where interface {
	Id() string
	String() string
	QuoteValue(value string) string
}

type WhereFactory interface {
	NewEmpty() Where
	New(args []string) Where
	NewQuery(field string, op string, builder Builder) Where
	NewArray(field string, op string, array []string) Where
}

type QuoteValue func(string) string

type BaseWhere struct {
	Field string
	Op    string
	Value string
	Query Builder
	Array []string
	Qv    QuoteValue

	id string
}

func (where *BaseWhere) Id() string {
	if where.id == "" {
		where.id = helpers.RandString(32)
	}
	return where.id
}

func (where *BaseWhere) String() string {
	value := ""
	switch {
	case where.Query != nil:
		value = "(" + where.Query.Query() + ")"
	case where.Value != "":
		value = where.Qv(where.Value)
	case where.Array != nil:
		value = "("
		length := len(where.Array)
		for i, item := range where.Array {
			value += where.Qv(item)
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

func (where *BaseWhere) QuoteValue(value string) string {
	return where.Qv(value)
}
