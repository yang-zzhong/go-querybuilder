package mysql

import . "yang-zzhong/database/query"

type MysqlBuilder struct {
	BaseBuilder
}

func New() Builder {
	baseBuilder := new(BaseBuilder)
	InitBuilder(baseBuilder, new(MysqlWhereFactory))

	builder := new(MysqlBuilder)
	builder.BaseBuilder = *baseBuilder

	return builder
}
