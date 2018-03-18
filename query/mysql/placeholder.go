package mysql

type MysqlPlaceholder struct{}

func (ph *MysqlPlaceholder) Ph() string {
	return "?"
}
