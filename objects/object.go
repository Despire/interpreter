package objects

import "fmt"

const (
	INTEGER Type = "INTEGER"
	BOOLEAN      = "BOOLEAN"
	NULL         = "NULL"
)

type (
	Type string

	Object interface {
		Type() Type
		Inspect() string
	}
)

type (
	Integer struct {
		Value int64
	}

	Boolean struct {
		Value bool
	}

	Null struct{}
)

// implement Object interface
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() Type      { return INTEGER }

// implement Object interface
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() Type      { return BOOLEAN }

// implement Object interface
func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() Type      { return NULL }
