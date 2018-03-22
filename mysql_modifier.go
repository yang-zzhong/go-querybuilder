package querybuilder

import "regexp"

type MysqlModifier struct {
	BaseModifier
}

func (mm *MysqlModifier) Ph(_ string) string {
	return "?"
}

func (mm *MysqlModifier) QuoteName(name string) string {
	point := regexp.MustCompile("\\.")
	return "`" + (string)(point.ReplaceAll(([]byte)(name), ([]byte)("`.`"))) + "`"
}
