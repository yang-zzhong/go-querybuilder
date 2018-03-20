package oracle

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
