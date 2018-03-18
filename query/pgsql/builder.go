package mysql

import . "yang-zzhong/database/query"

type PgsqlBuilder struct {
	BaseBuilder
}

func New() Builder {
	baseBuilder := new(BaseBuilder)
	whereFactory := new(BaseWhereFactory)
	whereFactory.Ph = new(PgsqlPlaceholder)
	InitBuilder(baseBuilder, whereFactory, whereFactory.Ph)

	builder := new(PgsqlBuilder)
	builder.BaseBuilder = *baseBuilder

	return builder
}
