package query

type Placeholder interface {
	Ph(name string) string
	PrePh() string
	PhRegExp() string
}
