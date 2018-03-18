package mysql

import (
	"yang-zzhong/helpers"
)

type PgsqlPlaceholder struct{}

func (ph *PgsqlPlaceholder) Ph() string {
	return "$" + helpers.RandString(10)
}
