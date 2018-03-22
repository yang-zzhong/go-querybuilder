package querybuilder

import (
	helpers "github.com/yang-zzhong/go-helpers"
	"regexp"
	"strconv"
)

//
// Quote's callback type
//
type QuoteWhere func(builder *Builder)

const (
	DESC string = "DESC"
	ASC  string = "ASC"
)

const (
	t_WHERE = iota
	t_RAW
	t_QUOTE_BEGIN
	t_QUOTE_END
	t_AND
	t_OR
)

type Condition struct {
	t  int    // condition type
	id string // condition id
}

type Builder struct {
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

	whereFactory *WhereFactory

	values []interface{}

	limit int

	offset int

	modifier Modifier

	replace bool
}

func NewBuilder(modifier Modifier) *Builder {
	builder := new(Builder)
	builder.conditions = []Condition{}
	builder.wheres = make(map[string]Where)
	builder.selects = []string{"*"}
	builder.orders = make(map[string]string)
	builder.whereFactory = NewWF(modifier)
	builder.values = []interface{}{}
	builder.limit = -1
	builder.offset = -1
	builder.modifier = modifier
	builder.replace = true

	return builder
}

func (builder *Builder) Replace(replace bool) *Builder {
	builder.replace = replace
	return builder
}

func (builder *Builder) From(table string) *Builder {
	builder.table = table
	return builder
}

func (builder *Builder) Select(cols []string) *Builder {
	builder.selects = cols
	return builder
}

func (builder *Builder) OrderBy(field string, order string) *Builder {
	builder.orders[field] = order
	return builder
}

func (builder *Builder) Where(args ...string) *Builder {
	where := builder.makeWhere(args)
	builder.conditions = append(builder.conditions, Condition{t_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *Builder) WhereRaw(raw string) *Builder {
	condition := Condition{t_RAW, raw}
	builder.conditions = append(builder.conditions, condition)

	return builder
}

func (builder *Builder) Quote(call QuoteWhere) *Builder {
	builder.conditions = append(builder.conditions, Condition{t_QUOTE_BEGIN, ""})
	call(builder)
	builder.conditions = append(builder.conditions, Condition{t_QUOTE_END, ""})

	return builder
}

func (builder *Builder) WhereInQuery(field string, ins *Builder) *Builder {
	return builder.WhereQuery(field, IN, ins)
}

func (builder *Builder) WhereNotInQuery(field string, ins *Builder) *Builder {
	return builder.WhereQuery(field, NOTIN, ins)
}

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

func (builder *Builder) WhereIn(field string, ins []string) *Builder {
	where := builder.makeArrayWhere(field, IN, ins)
	builder.conditions = append(builder.conditions, Condition{t_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *Builder) WhereNotIn(field string, ins []string) *Builder {
	where := builder.makeArrayWhere(field, NOTIN, ins)
	builder.conditions = append(builder.conditions, Condition{t_WHERE, where.Id()})
	builder.wheres[where.Id()] = where

	return builder
}

func (builder *Builder) Limit(limit int) *Builder {
	builder.limit = limit
	return builder
}

func (builder *Builder) Offset(offset int) *Builder {
	builder.offset = offset
	return builder
}

func (builder *Builder) Params() []interface{} {
	result := []interface{}{}
	for _, val := range builder.values {
		result = append(result, val)
	}

	return result
}

func (builder *Builder) ForQuery() string {
	builder.values = []interface{}{}
	selects := builder.quoteSelects()
	sql := "SELECT " + helpers.Implode(selects, ", ") + " FROM " + builder.modifier.QuoteName(builder.table)
	if len(builder.conditions) > 0 {
		sql += " WHERE " + builder.handleWhere()
	}
	if len(builder.orders) > 0 {
		sql += handleOrderBy(builder.orders)
	}
	if builder.limit > -1 {
		sql += " LIMIT " + strconv.Itoa(builder.limit)
	}
	if builder.offset > -1 {
		sql += " OFFSET " + strconv.Itoa(builder.offset)
	}
	if builder.replace {
		return replace(builder.modifier, sql)
	}

	return sql
}

func (builder *Builder) quoteTable() string {
	return builder.modifier.QuoteName(builder.table)
}

func (builder *Builder) quoteSelects() []string {
	result := make([]string, len(builder.selects))
	for i, item := range builder.selects {
		result[i] = builder.modifier.QuoteName(item)
	} 

	return result
}

func (builder *Builder) ForRemove() string {
	builder.values = []interface{}{}
	sql := "DELETE FROM " + builder.table
	if len(builder.conditions) > 0 {
		sql += " WHERE " + builder.handleWhere()
	}
	if builder.limit > -1 {
		sql += " LIMIT " + strconv.Itoa(builder.limit)
	}
	if builder.offset > -1 {
		sql += " OFFSET " + strconv.Itoa(builder.offset)
	}

	return replace(builder.modifier, sql)
}

func (builder *Builder) ForUpdate(data map[string]string) string {
	builder.values = []interface{}{}
	sql := "UPDATE " + builder.table + " SET "
	length := len(data)
	i := 1
	for field, value := range data {
		sql += field + "=" + builder.modifier.PrePh()
		builder.values = append(builder.values, value)
		if i < length {
			sql += ", "
		}
		i++
	}
	if len(builder.conditions) > 0 {
		sql += " WHERE " + builder.handleWhere()
	}
	if builder.limit > -1 {
		sql += " LIMIT " + strconv.Itoa(builder.limit)
	}
	if builder.offset > -1 {
		sql += " OFFSET " + strconv.Itoa(builder.offset)
	}

	return replace(builder.modifier, sql)
}

func (builder *Builder) makeWhere(args []string) Where {
	return builder.whereFactory.New(args)
}

func (builder *Builder) makeQueryWhere(field string, op string, other *Builder) Where {
	return builder.whereFactory.NewQuery(field, op, other)
}

func (builder *Builder) makeArrayWhere(field string, op string, array []string) Where {
	return builder.whereFactory.NewArray(field, op, array)
}

func (builder *Builder) handleWhere() string {
	wheres := []string{}
	for _, condi := range builder.conditions {
		length := len(wheres)
		if length == 0 {
			where := builder.wheres[condi.id]
			wheres = append(wheres, where.String())
			for _, value := range where.Params() {
				builder.values = append(builder.values, value)
			}
			continue
		}
		last := wheres[length-1]
		if condi.t == t_WHERE {
			where := builder.wheres[condi.id]
			wheres = addAnd(wheres, last)
			wheres = append(wheres, where.String())
			for _, value := range where.Params() {
				builder.values = append(builder.values, value)
			}
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

func replace(modifier Modifier, src string) string {
	bSrc := ([]byte)(src)
	search := regexp.MustCompile(modifier.PrePh())
	i := 1
	result := search.ReplaceAllFunc(bSrc, func(matched []byte) []byte {
		res := modifier.Ph(strconv.Itoa(i))
		i++
		return ([]byte)(res)
	})

	return (string)(result)
}
