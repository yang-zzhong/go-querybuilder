package querybuilder

import (
	"regexp"
	"strconv"
)

func (builder *Builder) Replace(replace bool) *Builder {
	builder.replace = replace
	return builder
}

func replace(modifier Modifier, src string) string {
	bSrc := ([]byte)(src)
	search := regexp.MustCompile(modifier.PrePh())
	i := 1
	result := search.ReplaceAllFunc(bSrc, func(matched []byte) []byte {
		res := modifier.Ph(strconv.Itoa(i))
		i++
		return ([]byte)(res)
	})

	return (string)(result)
}
