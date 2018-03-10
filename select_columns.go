package querybuilder

//
// select columns
// builder.Select("name", "age", "id")
//
func (builder *Builder) Select(cols ...string) *Builder {
	builder.selects = cols
	return builder
}
