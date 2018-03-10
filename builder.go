package querybuilder

import (
	"strconv"
	str "strings"
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

func (builder *Builder) From(table string) *Builder {
	builder.table = table
	return builder
}

func (builder *Builder) OrderBy(field string, order string) *Builder {
	builder.orders[field] = order
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
	columns := builder.QuotedColumns()
	sql := "SELECT " + str.Join(columns, ", ") + " FROM " + builder.QuotedTableName()
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

func (builder *Builder) QuotedTableName() string {
	return builder.modifier.QuoteName(builder.table)
}

func (builder *Builder) QuotedColumns() []string {
	result := make([]string, len(builder.selects))
	for i, item := range builder.selects {
		result[i] = builder.modifier.QuoteName(item)
	}

	return result
}

func (builder *Builder) ForRemove() string {
	builder.values = []interface{}{}
	sql := "DELETE FROM " + builder.QuotedTableName()
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
	sql := "UPDATE " + builder.QuotedTableName() + " SET "
	length := len(data)
	i := 1
	for field, value := range data {
		sql += builder.modifier.QuoteName(field) + "=" + builder.modifier.PrePh()
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
