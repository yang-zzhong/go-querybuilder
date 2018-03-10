package querybuilder

type Counter interface {
	ForCount() string
}

func (builder *Builder) ForCount() string {
	builder.Select("COUNT(1)")
	return builder.ForQuery()
}
