package main

import (
	"fmt"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() Environment {
	return Environment{
		enclosing: nil,
		values:    make(map[string]any),
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Environment) getAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Environment) assignAt(distance int, name Token, value any) {
	e.ancestor(distance).values[name.lexeme] = value
}

func (e *Environment) get(name Token) (any, error) {
	v, ok := e.values[name.lexeme]
	if ok {
		return v, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, RuntimeError{name, fmt.Sprintf("Undefined variable %q.", name.lexeme)}
}
