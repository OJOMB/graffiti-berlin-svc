package nanoID

import gonanoid "github.com/matoous/go-nanoid"

type Generator struct {
	symbols string
	length  int
}

func NewGenerator(symbols string, length int) *Generator {
	return &Generator{symbols: symbols, length: length}
}

func (gen *Generator) New() (string, error) {
	return gonanoid.Generate(gen.symbols, gen.length)
}
