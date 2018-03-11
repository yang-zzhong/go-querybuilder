package querybuilder

import (
	"strings"
)

func (builder *Builder) GroupBy(fieldNames ...string) *Builder {
	groupFields := []string{}
	for _, fieldName := range fieldNames {
		groupFields = append(groupFields, builder.modifier.QuoteName(fieldName))
	}

	builder.group = "GROUP BY " + strings.Join(groupFields, ", ")
	return builder
}

func (builder *Builder) Group(group string) *Builder {
	builder.group = group
	return builder
}
