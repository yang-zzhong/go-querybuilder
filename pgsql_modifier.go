package querybuilder

import "regexp"

type PgsqlModifier struct {
	BaseModifier
}

func (pm *PgsqlModifier) Ph(name string) string {
	return "$" + name
}

func (pm *PgsqlModifier) QuoteName(name string) string {
	point := regexp.MustCompile("\\.")
	return "\"" + (string)(point.ReplaceAll(([]byte)(name), ([]byte)("\".\""))) + "\""
}
