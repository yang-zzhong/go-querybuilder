package query

import (
	"regexp"
	"testing"
)

func TestBaseWhere(t *testing.T) {
	where := new(BaseWhere)
	InitBaseWhere(where)
	result := make(map[string]string)
	result[EQ] = "^field = @.*$"
	result[NEQ] = "^field != @.*$"
	result[GT] = "^field > @.*$"
	result[GTE] = "^field >= @.*$"
	result[LT] = "^field < @.*$"
	result[LTE] = "^field <= @.*$"
	result[LIKE] = "^field LIKE @.*$"
	result[NULL] = "^field IS NULL$"
	result[NOTNULL] = "^field IS NOT NULL$"
	result[IN] = "^field IN (.*)$"
	result[NOTIN] = "^field NOT IN (.*)$"
	where.Field = "field"
	for op, r := range result {
		where.Op = op
		switch op {
		case LIKE:
			where.Value = "%value%"
		case NULL:
		case NOTNULL:
		case IN:
			where.Array = []string{"value1", "value2"}
		case NOTIN:
			where.Array = []string{"value1", "value2"}
		default:
			where.Value = "value"
		}
		suc, _ := regexp.Match(r, ([]byte)(where.String()))
		if !suc {
			t.Error("test Error at ", op)
		}
	}
}
