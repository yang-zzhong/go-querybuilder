package query

type WhereFactory interface {
	New(args []string) Where
	NewQuery(field string, op string, builder Builder) Where
	NewArray(field string, op string, array []string) Where
}

type BaseWhereFactory struct {
	Ph Placeholder
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
	where := new(BaseWhere)
	InitBaseWhere(where, factory.Ph)
	where.Field = condi[0]
	where.Op = condi[1]
	where.Value = condi[2]

	return where
}

func (factory *BaseWhereFactory) NewQuery(field string, op string, other Builder) Where {
	where := new(BaseWhere)
	InitBaseWhere(where, factory.Ph)
	where.Field = field
	where.Op = op
	where.Query = other

	return where
}

func (factory *BaseWhereFactory) NewArray(field string, op string, array []string) Where {
	where := new(BaseWhere)
	InitBaseWhere(where, factory.Ph)
	where.Field = field
	where.Op = op
	where.Array = array

	return where
}
