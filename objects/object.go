package objects

import "fmt"

const (
	INTEGER Type = "INTEGER"
	BOOLEAN      = "BOOLEAN"
	NULL         = "NULL"
	RETURN       = "RETURN_VALUE"
	ERROR        = "ERROR"
)

type (
	Type string

	Object interface {
		Type() Type
		Inspect() string
	}
)

type (
	Return struct {
		Value Object
	}
	Integer struct {
		Value int64
	}

	Boolean struct {
		Value bool
	}

	Null struct{}

	Error struct {
		Value string
	}
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

// implement Object interface
func (r *Return) Inspect() string { return r.Value.Inspect() }
func (r *Return) Type() Type      { return RETURN }

// implement Object interface
func (e *Error) Inspect() string { return "ERROR: " + e.Value }
func (e *Error) Type() Type      { return ERROR }
