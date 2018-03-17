package query

import (
	helpers "github.com/yang-zzhong/go-helpers"
)

//
// QWhere's callback type
//
type QuoteWhere func(builder Builder)

type Builder interface {

	// Select Update OR Delete From whitch table
	From(table_name string) Builder

	// select columns to fetch, default ["*"]
	Select(cols []string) Builder

	// set query condition
	// query.Where(field, value)
	// query.Where(field, op, value)
	// # sample
	//		query.Where("age", 12).Where("name", "like", "%young").Where("age", ">", 12)
	Where(args ...string) Builder

	// quote a group of condition
	// # sample
	// 	 query.QWhere(func(query Builder) {
	// 		query.Where("hello", "world")
	// 		query.Or().Where("hello", "dlrow")
	//   })
	// will generate `(hello = "world" OR hello = "dlrow")`
	QWhere(call QuoteWhere) Builder

	// set a raw condition
	// # sample
	//	query.WhereRaw("name = \"world\"")
	WhereRaw(condition string) Builder

	// set in condition, the value is another builder, the builder must select a field
	// #sample
	//	anotherQuery.Select("id")
	//	query.WhereInQuery("id", anotherQuery)
	WhereInQuery(field string, query Builder) Builder

	// same as WhereInQuery
	WhereNotInQuery(field string, query Builder) Builder

	// set in condition, the vaue is a slice
	// #sample
	//	query.WhereIn("name", []string{"狗蛋", "二狗子"})
	WhereIn(field string, ins []string) Builder

	// same as WhereIn
	WhereNotIn(field string, ins []string) Builder

	// or connected two conditions
	Or() Builder

	// and connected two conditions
	And() Builder

	// order the result
	// #sample
	//	query.OrderBy("name", DESC)
	//	query.OrderBy("age", ASC)
	OrderBy(filed string, order string) Builder

	// execute select query
	Query() string

	// execute update query
	Update(map[string]string) string

	// execute remove query
	Remove() string
}

const (
	DESC string = "DESC"
	ASC  string = "ASC"
)

const (
	T_WHERE = iota
	T_RAW
	T_QUOTE_BEGIN
	T_QUOTE_END
	T_AND
	T_OR
)

type Condition struct {
	t  int    // condition type
	id string // condition id
}

type BaseBuilder struct {
	// users
	table string

	conditions []Condition

	wheres map[string]Where

	// ["name", "age"]
	// ["*"]
	selects []string

	//
	// ["age" => "desc"]
	//
	orders map[string]string

	whereFactory WhereFactory
}

func InitBuilder(builder *BaseBuilder, where WhereFactory) {
	builder.conditions = []Condition{}
	builder.wheres = make(map[string]Where)
	builder.selects = []string{}
	builder.orders = make(map[string]string)
	builder.whereFactory = where
}

func (builder *BaseBuilder) From(table string) Builder {
	builder.table = table

	return builder
}

func (builder *BaseBuilder) Select(cols []string) Builder {
	builder.selects = cols

	return builder
}

func (builder *BaseBuilder) OrderBy(field string, order string) Builder {
	builder.orders[field] = order
	return builder
}

func (builder *BaseBuilder) Where(args ...string) Builder {
	where := builder.makeWhere(args)
	builder.conditions = append(builder.conditions, Condition{T_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *BaseBuilder) WhereRaw(raw string) Builder {
	condition := Condition{T_RAW, raw}
	builder.conditions = append(builder.conditions, condition)

	return builder
}

func (builder *BaseBuilder) QWhere(call QuoteWhere) Builder {
	builder.conditions = append(builder.conditions, Condition{T_QUOTE_BEGIN, ""})
	call(builder)
	builder.conditions = append(builder.conditions, Condition{T_QUOTE_END, ""})

	return builder
}

func (builder *BaseBuilder) WhereInQuery(field string, ins Builder) Builder {
	return builder.WhereQuery(field, IN, ins)
}

func (builder *BaseBuilder) WhereNotInQuery(field string, ins Builder) Builder {
	return builder.WhereQuery(field, NOTIN, ins)
}

func (builder *BaseBuilder) WhereQuery(field string, op string, other Builder) Builder {
	where := builder.makeQueryWhere(field, op, other)
	builder.conditions = append(builder.conditions, Condition{T_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *BaseBuilder) And() Builder {
	builder.conditions = append(builder.conditions, Condition{T_QUOTE_END, ""})

	return builder
}

func (builder *BaseBuilder) Or() Builder {
	builder.conditions = append(builder.conditions, Condition{T_OR, ""})

	return builder
}

func (builder *BaseBuilder) WhereIn(field string, ins []string) Builder {
	where := builder.makeArrayWhere(field, IN, ins)
	builder.conditions = append(builder.conditions, Condition{T_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *BaseBuilder) WhereNotIn(field string, ins []string) Builder {
	where := builder.makeArrayWhere(field, NOTIN, ins)
	builder.conditions = append(builder.conditions, Condition{T_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *BaseBuilder) Query() string {

	selects := builder.selects
	if selects == nil {
		selects = []string{"*"}
	}
	sql := "SELECT " + helpers.Implode(selects, ", ") + " FROM " + builder.table
	if builder.conditions != nil {
		sql += " WHERE " + builder.handleWhere()
	}
	if builder.orders != nil {
		sql += handleOrderBy(builder.orders)
	}

	return sql
}

func (builder *BaseBuilder) Remove() string {
	sql := "DELETE FROM " + builder.table
	if builder.conditions != nil {
		sql += " WHERE " + builder.handleWhere()
	}

	return sql
}

func (builder *BaseBuilder) Update(data map[string]string) string {
	return ""
}

func (builder *BaseBuilder) makeWhere(args []string) Where {
	return builder.whereFactory.New(args)
}

func (builder *BaseBuilder) makeQueryWhere(field string, op string, other Builder) Where {
	return builder.whereFactory.NewQuery(field, op, other)
}

func (builder *BaseBuilder) makeArrayWhere(field string, op string, array []string) Where {
	return builder.whereFactory.NewArray(field, op, array)
}

func (builder *BaseBuilder) handleWhere() string {
	wheres := []string{}
	for _, condi := range builder.conditions {
		length := len(wheres)
		if length == 0 {
			wheres = append(wheres, builder.wheres[condi.id].String())
			continue
		}
		last := wheres[len(wheres)-1]
		if condi.t == T_WHERE {
			wheres = addAnd(wheres, last)
			wheres = append(wheres, builder.wheres[condi.id].String())
			continue
		}
		where := ""
		switch condi.t {
		case T_QUOTE_BEGIN:
			wheres = addAnd(wheres, last)
			where = "("
		case T_QUOTE_END:
			where = ")"
		case T_AND:
			where = " AND "
		case T_OR:
			where = " OR "
		case T_RAW:
			where = condi.id
		}
		wheres = append(wheres, where)
	}

	return helpers.Implode(wheres, "")
}

func handleOrderBy(orders map[string]string) string {
	sql := " ORDER BY "
	length := len(orders)
	i := 1
	for field, order := range orders {
		sql += field + " " + order
		if i < length {
			sql += ", "
		}
		i++
	}

	return sql
}

func addAnd(wheres []string, last string) []string {
	if last == "(" {
		return wheres
	}
	if last != " AND " && last != " OR " {
		wheres = append(wheres, " AND ")
	}

	return wheres
}
