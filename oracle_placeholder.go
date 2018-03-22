package querybuilder

import "regexp"

type OracleModifier struct {
	BaseModifier
}

func (om *OracleModifier) Ph(name string) string {
	return ":" + name
}

func (om *OracleModifier) QuoteName(name string) string {
	point := regexp.MustCompile("\\.")
	return "\"" + (string)(point.ReplaceAll(([]byte)(name), ([]byte)("\".\""))) + "\""
}
