package main

import (
	"fmt"
)

type Environment struct {
	values map[string]any
}

func NewEnvironment() Environment {
	return Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name Token) (any, error) {
	v, ok := e.values[name.lexeme]
	if ok {
		return v, nil
	}

	return nil, RuntimeError{name, fmt.Sprintf("Undefined variable %q.", name.lexeme)}
}
