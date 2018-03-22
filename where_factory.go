package querybuilder

type WhereFactory interface {
	New(args []string) Where
	NewQuery(field string, op string, builder Builder) Where
	NewArray(field string, op string, array []string) Where
}

type BaseWhereFactory struct {
	ph Placeholder
}

func NewWF(ph Placeholder) WhereFactory {
	wf := new(BaseWhereFactory)
	wf.ph = ph

	return wf
}

func (factory *BaseWhereFactory) New(args []string) Where {
	length := len(args)
	condi := []string{}
	switch length {
	case 2:
		condi = []string{args[0], "=", args[1]}
	case 3:
		condi = args
	}
	where := NewW(factory.ph)
	where.Field = condi[0]
	where.Op = condi[1]
	where.Value = condi[2]

	return where
}

func (factory *BaseWhereFactory) NewQuery(field string, op string, other Builder) Where {
	where := NewW(factory.ph)
	where.Field = field
	where.Op = op
	where.Query = other

	return where
}

func (factory *BaseWhereFactory) NewArray(field string, op string, array []string) Where {
	where := NewW(factory.ph)
	where.Field = field
	where.Op = op
	where.Array = array

	return where
}
