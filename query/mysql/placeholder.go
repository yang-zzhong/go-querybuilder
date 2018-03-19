package mysql

type MysqlPlaceholder struct{}

func (ph *MysqlPlaceholder) Ph(name string) string {
	name = "?"
	return name
}

func (ph *MysqlPlaceholder) PrePh() string {
	return "__@__"
}
