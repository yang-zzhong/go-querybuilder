package sqlite

import . "godb/query"

type SqliteBuilder struct {
	BaseBuilder
}

func New() Builder {
	baseBuilder := new(BaseBuilder)
	whereFactory := new(BaseWhereFactory)
	whereFactory.Ph = new(SqlitePlaceholder)
	InitBuilder(baseBuilder, whereFactory, whereFactory.Ph)

	builder := new(SqliteBuilder)
	builder.BaseBuilder = *baseBuilder

	return builder
}
