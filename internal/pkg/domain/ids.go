package domain

type IDTool interface {
	IDGenerator
	IDValidator
}

type IDGenerator interface {
	New() (string, error)
}

type IDValidator interface {
	IsValid(ID string) bool
}
