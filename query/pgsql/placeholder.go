package pgsql

type PgsqlPlaceholder struct{}

func (ph *PgsqlPlaceholder) Ph(name string) string {
	return "$" + name
}

func (ph *PgsqlPlaceholder) PrePh() string {
	return "__@__"
}
