package querybuilder

func (builder *Builder) Page(page int, pageSize int) *Builder {
	return builder.Offset(pageSize * (page - 1)).Limit(pageSize)
}
