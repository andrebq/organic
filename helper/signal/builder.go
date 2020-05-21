package signal

import "github.com/andrebq/organic/cell"

type (
	// Builder provides a type-safe interface to create signals
	Builder struct {
		items []interface{}
	}
)

// Build creates an empty signal builder
func Build() *Builder {
	return &Builder{}
}

// Str appends the given string to the signal
func (b *Builder) Str(s string) *Builder {
	b.items = append(b.items, s)
	return b
}

// Bytes appends a sequence of bytes to the builder
func (b *Builder) Bytes(v []byte) *Builder {
	b.items = append(b.items, v)
	return b
}

// Signal returns a valid signal with all the data that was encoded
func (b *Builder) Signal() cell.Signal {
	s, err := cell.EncodeSignal(b.items...)
	if err != nil {
		panic("should never happen: " + err.Error())
	}
	return s
}
