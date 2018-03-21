package model

import (
	"database/sql"
)

type EachCall func(Model)
type Collection struct {
	m    Model
	rows Rows
}

func NewCollection(rows Rows, m Model) {
	col := new(Collection)
	col.m = m
	col.rows = rows

	return col
}

func (col *Collection) First() Model {

}

func (col *Collection) AddModel(m Model) {

}

func (col *Collection) AddRow(row Row) {

}

func (col *Collection) Each(call EachCall) {

}
