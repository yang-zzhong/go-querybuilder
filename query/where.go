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
	Params() []interface{}
}

type BaseWhere struct {
	Field string
	Op    string
	Value string
	Query Builder
	Array []string

	id          string
	values      []interface{}
	placeholder Placeholder
}

func InitBaseWhere(where *BaseWhere, ph Placeholder) {
	where.values = []interface{}{}
	where.id = helpers.RandString(32)
	where.placeholder = ph
}

// func (where *BaseWhere) Params() map[string]string {
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
		value = "(" + where.Query.ForQuery() + ")"
		where.values = where.Query.Params()
	case where.Value != "":
		value = where.placeholder.PrePh()
		where.values = []interface{}{where.Value}
	case where.Array != nil:
		value = "("
		length := len(where.Array)
		for i, item := range where.Array {
			value += where.placeholder.PrePh()
			where.values = append(where.values, item)
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
