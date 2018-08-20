package querybuilder

import (
	str "strings"
)

type QueryWhere interface {
	Where(...string) *Builder
	WhereRaw(string) *Builder
	WhereIn(string, []interface{}) *Builder
	WhereNotIn(string, []string) *Builder
	WhereQuery(string, string, *Builder) *Builder
	WhereInQuery(string, *Builder) *Builder
	WhereNotInQuery(string, *Builder) *Builder
	Quote(QuoteWhere) *Builder
	And() *Builder
	Or() *Builder
}

//
// config the query conditions, it takes three params
// param field_name for the first param
// op/value for the second param, if the param is value, the op
// is EQ
// if the second param is op, then the third param is the value
// all params type is string, this means if has a int type value
// you must quote it as a string
//
// builder.Where("age", GT, "15").
//	Where("name", LIKE "Young%").
//	Where("id", 1)
//
func (builder *Builder) Where(args ...interface{}) *Builder {
	where := builder.makeWhere(args)
	builder.conditions = append(builder.conditions, Condition{t_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

//
// some time we need a complicated condition, you can write a where raw string
// builder.WhereRaw("cities && {'重庆', '北京', '上海'}")
// PLUS: 重庆, 北京, 上海 are all Chinese City, first of them(Chongqing) that I living in
//
func (builder *Builder) WhereRaw(raw string) *Builder {
	condition := Condition{t_RAW, raw}
	builder.conditions = append(builder.conditions, condition)

	return builder
}

//
// some times we need quote a group of conditions to assign a solid priority, we can use Quote
// builder.Where("age", GT, "15").Quote(func (builder *Builder) {
//		builder.Where("name", LIKE, "%Young").Or().Where("name", LIKE, "%Old")
// })
// will generate: age > 15 AND (name LIKE '%Young' OR name LIKE '%Old')
//
func (builder *Builder) Quote(call QuoteWhere) *Builder {
	builder.conditions = append(builder.conditions, Condition{t_QUOTE_BEGIN, ""})
	call(builder)
	builder.conditions = append(builder.conditions, Condition{t_QUOTE_END, ""})

	return builder
}

//
// see WhereQuery
//
func (builder *Builder) WhereInQuery(field string, ins *Builder) *Builder {
	return builder.WhereQuery(field, IN, ins)
}

//
// see WhereQuery
//
func (builder *Builder) WhereNotInQuery(field string, ins *Builder) *Builder {
	return builder.WhereQuery(field, NOTIN, ins)
}

//
// where query takes a sub query as a condition
// authors := builder.Select("author_id").WhereIn("category", []string{"Golang", "C++"})
// users := builder.WhereQuery("user_id", IN, authors) -----eq---- users := builder.WhereInQuery("user_id", authors)
//
func (builder *Builder) WhereQuery(field string, op string, other *Builder) *Builder {
	where := builder.makeQueryWhere(field, op, other)
	builder.conditions = append(builder.conditions, Condition{t_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *Builder) And() *Builder {
	builder.conditions = append(builder.conditions, Condition{t_QUOTE_END, ""})

	return builder
}

func (builder *Builder) Or() *Builder {
	builder.conditions = append(builder.conditions, Condition{t_OR, ""})

	return builder
}

func (builder *Builder) WhereIn(field string, ins []interface{}) *Builder {
	where := builder.makeArrayWhere(field, IN, ins)
	builder.conditions = append(builder.conditions, Condition{t_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *Builder) WhereNotIn(field string, ins []interface{}) *Builder {
	where := builder.makeArrayWhere(field, NOTIN, ins)
	builder.conditions = append(builder.conditions, Condition{t_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *Builder) makeWhere(args []interface{}) Where {
	return builder.whereFactory.New(args)
}

func (builder *Builder) makeQueryWhere(field string, op string, other *Builder) Where {
	return builder.whereFactory.NewQuery(field, op, other)
}

func (builder *Builder) makeArrayWhere(field string, op string, array []interface{}) Where {
	return builder.whereFactory.NewArray(field, op, array)
}

func (builder *Builder) handleWhere() string {
	wheres := []string{}
	for _, condi := range builder.conditions {
		length := len(wheres)
		last := ""
		if length > 0 {
			last = wheres[length-1]
		}
		if condi.t == t_WHERE {
			where := builder.wheres[condi.id]
			wheres = addAnd(wheres, last)
			wheres = append(wheres, where.String())
			for _, value := range where.Params() {
				builder.values = append(builder.values, value)
			}
			continue
		}
		if condi.t == t_RAW {
			wheres = addAnd(wheres, last)
			wheres = append(wheres, condi.id)
			continue
		}
		where := ""
		switch condi.t {
		case t_QUOTE_BEGIN:
			wheres = addAnd(wheres, last)
			where = "("
		case t_QUOTE_END:
			where = ")"
		case t_AND:
			where = " AND "
		case t_OR:
			where = " OR "
		case t_RAW:
			where = condi.id
		}
		wheres = append(wheres, where)
	}

	return str.Join(wheres, "")
}

func addAnd(wheres []string, last string) []string {
	if last == "(" || last == "" {
		return wheres
	}
	if last != " AND " && last != " OR " {
		wheres = append(wheres, " AND ")
	}

	return wheres
}
