package model

type Model interface {
	TableName() string
	IdKey() interface{}
	Json() string
	Map() map[string]interface{}
}

func Find(m Model, id interface{}) Model {
}

func Queryer(m Model) Queryer {
	queryer := NewQueryer()

	return queryer.From(m.TableName())
}

func Create(m Model) {

}

func Save(m Model) {

}
