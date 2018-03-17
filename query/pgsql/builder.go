package mysql

import . "yang-zzhong/database/query"

type PgsqlBuilder struct {
	BaseBuilder
}

func New() Builder {
	baseBuilder := new(BaseBuilder)
	InitBuilder(baseBuilder, new(PgsqlWhereFactory))

	builder := new(PgsqlBuilder)
	builder.BaseBuilder = *baseBuilder

	return builder
}
