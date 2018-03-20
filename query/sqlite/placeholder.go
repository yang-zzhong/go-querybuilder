package sqlite

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
