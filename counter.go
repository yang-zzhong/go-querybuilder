package querybuilder

func (builder *Builder) ForCount() string {
	builder.Select(E{"COUNT(1) as count"})
	return builder.ForQuery()
}
