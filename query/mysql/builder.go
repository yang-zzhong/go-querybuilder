package mysql

import . "gdb/query"

type MysqlBuilder struct {
	BaseBuilder
}

func New() Builder {
	baseBuilder := new(BaseBuilder)
	whereFactory := new(BaseWhereFactory)
	whereFactory.Ph = new(MysqlPlaceholder)
	InitBuilder(baseBuilder, whereFactory, whereFactory.Ph)

	builder := new(MysqlBuilder)
	builder.BaseBuilder = *baseBuilder

	return builder
}
