package mysql

import (
	. "yang-zzhong/database/query"
)

type PgsqlWhere struct {
	BaseWhere
}

type PgsqlWhereFactory struct{}

func (factory *PgsqlWhereFactory) NewEmpty() Where {
	baseWhere := new(BaseWhere)
	baseWhere.Qv = func(value string) string {
		return "'" + value + "'"
	}
	where := new(PgsqlWhere)
	where.BaseWhere = *baseWhere

	return where
}

func (factory *PgsqlWhereFactory) New(args []string) Where {
	length := len(args)
	condi := []string{}
	switch length {
	case 2:
		condi = []string{args[0], "=", args[1]}
	case 3:
		condi = args
	}
	where := new(BaseWhere)
	where.Field = condi[0]
	where.Op = condi[1]
	where.Value = condi[2]
	where.Qv = func(value string) string {
		return "'" + value + "'"
	}
	pgsqlWhere := new(PgsqlWhere)
	pgsqlWhere.BaseWhere = *where

	return pgsqlWhere
}

func (factory *PgsqlWhereFactory) NewQuery(field string, op string, other Builder) Where {
	where := new(BaseWhere)
	where.Field = field
	where.Op = op
	where.Query = other
	where.Qv = func(value string) string {
		return "'" + value + "'"
	}
	pgsqlWhere := new(PgsqlWhere)
	pgsqlWhere.BaseWhere = *where

	return pgsqlWhere
}

func (factory *PgsqlWhereFactory) NewArray(field string, op string, array []string) Where {
	where := new(BaseWhere)
	where.Field = field
	where.Op = op
	where.Array = array
	where.Qv = func(value string) string {
		return "'" + value + "'"
	}
	pgsqlWhere := new(PgsqlWhere)
	pgsqlWhere.BaseWhere = *where

	return pgsqlWhere
}
