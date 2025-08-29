package main

import (
	"fmt"
)

type LoxInstance struct {
	klass  *LoxClass
	fields map[string]any
}

func NewLoxInstance(klass *LoxClass) *LoxInstance {
	return &LoxInstance{
		klass:  klass,
		fields: make(map[string]any),
	}
}

func (l *LoxInstance) get(name Token) (any, error) {
	if val, ok := l.fields[name.lexeme]; ok {
		return val, nil
	}

	method := l.klass.findMethod(name.lexeme)
	if method != nil {
		return method, nil
	}

	return nil, RuntimeError{name, fmt.Sprintf("Undefined property %q.", name.lexeme)}
}

func (l *LoxInstance) set(name Token, value any) {
	l.fields[name.lexeme] = value
}

func (l *LoxInstance) String() string {
	return l.klass.name + " instance"
}
