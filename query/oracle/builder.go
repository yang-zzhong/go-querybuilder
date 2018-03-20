package oracle

import . "godb/query"

type OracleBuilder struct {
	BaseBuilder
}

func New() Builder {
	baseBuilder := new(BaseBuilder)
	whereFactory := new(BaseWhereFactory)
	whereFactory.Ph = new(OraclePlaceholder)
	InitBuilder(baseBuilder, whereFactory, whereFactory.Ph)

	builder := new(OracleBuilder)
	builder.BaseBuilder = *baseBuilder

	return builder
}
