package objects

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	o, ok := e.store[name]
	if !ok && e.outer != nil {
		o, ok = e.outer.Get(name)
	}
	return o, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
