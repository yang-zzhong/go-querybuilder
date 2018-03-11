package querybuilder

type Modifier interface {
	Ph(name string) string
	PrePh() string
	QuoteName(string) string
}

type BaseModifier struct{}

func (bm *BaseModifier) PrePh() string {
	return "__@__"
}
