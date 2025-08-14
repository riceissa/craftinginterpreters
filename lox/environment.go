package main

import (
	"fmt"
)

type Environment struct {
	values map[string]any
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name Token) (any, error) {
	v, ok := e.values[name.lexeme]
	if ok {
		return v, nil
	}

	// TODO: maybe use runtimeError here instead?
	return nil, fmt.Errorf("Undefined variable %q.", name.lexeme)
}
