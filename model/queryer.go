package model

import (
	. "github.com/yang-zzhong/godb"
)

type Queryer struct {
	model Model
	conn  interface{}
	BaseBuilder
}

func NewQueryer(m Model, db interface{}) {
	if db != nil {
		conn = db
	}
	queryer := new(Queryer)
	queryer.model = m
	queryer.conn = conn

	return queryer
}

func (queryer *Queryer) Query() Collection {
	rows, err := queryer.conn.Query(queryer.ForQuery(), queryer.Params()...)
	return NewCollection(rows, queryer.model)
}

func (queryer *Queryer) Update(data map[string]interface{}) {
	queryer.conn.Exec(queryer.ForUpdate(data), queryer.Params()...)
}

func (queryer *Queryer) Remove() Collection {
	queryer.conn.Exec(queryer.ForRemove(), queryer.Params()...)
}
