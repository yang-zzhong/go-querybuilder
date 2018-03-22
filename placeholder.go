package querybuilder

type Placeholder interface {
	Ph(name string) string
	PrePh() string
	PhRegExp() string
}

/////////////////mysql placeholder//////////////////
type MysqlPlaceholder struct{}

func (ph *MysqlPlaceholder) Ph(name string) string {
	name = "?"
	return name
}

func (ph *MysqlPlaceholder) PhRegExp() string {
	return "\\?"
}

func (ph *MysqlPlaceholder) PrePh() string {
	return "__@__"
}

///////////////pgsql placeholder/////////////////
type PgsqlPlaceholder struct{}

func (ph *PgsqlPlaceholder) Ph(name string) string {
	return "$" + name
}

func (ph *PgsqlPlaceholder) PhRegExp() string {
	return "\\$\\d+"
}

func (ph *PgsqlPlaceholder) PrePh() string {
	return "__@__"
}

/////////////oracle placeholder///////////////////
type OraclePlaceholder struct{}

func (ph *OraclePlaceholder) Ph(name string) string {
	return ":" + name
}

func (ph *OraclePlaceholder) PhRegExp() string {
	return "\\:\\d+"
}

func (ph *OraclePlaceholder) PrePh() string {
	return "__@__"
}

/////////////sqlite placeholder////////////////
type SqlitePlaceholder struct{}

func (ph *SqlitePlaceholder) Ph(name string) string {
	name = "?"
	return name
}

func (ph *SqlitePlaceholder) PhRegExp() string {
	return "\\?"
}

func (ph *SqlitePlaceholder) PrePh() string {
	return "__@__"
}
