package mysql

import (
	. "yang-zzhong/database/query"
)

type MysqlWhere struct {
	BaseWhere
}

type MysqlWhereFactory struct{}

func (factory *MysqlWhereFactory) NewEmpty() Where {
	baseWhere := new(BaseWhere)
	baseWhere.Qv = func(value string) string {
		return "\"" + value + "\""
	}
	where := new(MysqlWhere)
	where.BaseWhere = *baseWhere

	return where
}

func (factory *MysqlWhereFactory) New(args []string) Where {
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
		return "\"" + value + "\""
	}
	mysqlWhere := new(MysqlWhere)
	mysqlWhere.BaseWhere = *where

	return mysqlWhere
}

func (factory *MysqlWhereFactory) NewQuery(field string, op string, other Builder) Where {
	where := new(BaseWhere)
	where.Field = field
	where.Op = op
	where.Query = other
	where.Qv = func(value string) string {
		return "\"" + value + "\""
	}
	mysqlWhere := new(MysqlWhere)
	mysqlWhere.BaseWhere = *where

	return mysqlWhere
}

func (factory *MysqlWhereFactory) NewArray(field string, op string, array []string) Where {
	where := new(BaseWhere)
	where.Field = field
	where.Op = op
	where.Array = array
	where.Qv = func(value string) string {
		return "\"" + value + "\""
	}
	mysqlWhere := new(MysqlWhere)
	mysqlWhere.BaseWhere = *where

	return mysqlWhere
}
