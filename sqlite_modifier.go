package querybuilder

import "regexp"

type SqliteModifier struct {
	BaseModifier
}

func (sm *SqliteModifier) Ph(name string) string {
	return "$" + name
}

func (sm *SqliteModifier) QuoteName(name string) string {
	point := regexp.MustCompile("\\.")
	return "\"" + (string)(point.ReplaceAll(([]byte)(name), ([]byte)("\".\""))) + "\""
}
