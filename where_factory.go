package querybuilder

type WhereFactory struct {
	modifier Modifier
}

func NewWF(modifier Modifier) *WhereFactory {
	return &WhereFactory{modifier}
}

func (factory *WhereFactory) New(args []interface{}) Where {
	length := len(args)
	condi := []interface{}{}
	switch length {
	case 2:
		condi = []interface{}{args[0], "=", args[1]}
	case 3:
		condi = args
	}
	where := NewW(factory.modifier)
	where.Field = condi[0].(string)
	where.Op = condi[1].(string)
	where.Value = condi[2]

	return where
}

func (factory *WhereFactory) NewQuery(field string, op string, other *Builder) Where {
	where := NewW(factory.modifier)
	where.Field = field
	where.Op = op
	where.Query = other

	return where
}

func (factory *WhereFactory) NewArray(field string, op string, array []string) Where {
	where := NewW(factory.modifier)
	where.Field = field
	where.Op = op
	where.Array = array

	return where
}
