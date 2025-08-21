package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration Function
	closure *Environment
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %v >", f.declaration.name.lexeme)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	environment := NewEnvironment()
	environment.enclosing = f.closure
	for i := range len(f.declaration.params) {
		environment.define(f.declaration.params[i].lexeme, arguments[i])
	}

	result, err := interpreter.executeBlock(f.declaration.body, &environment)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result.value, nil
	}
	return nil, nil
}

type LoxNativeFunction struct {
	arity int
	fn    func(*Interpreter, []any) any
	name  string
}

func (n *LoxNativeFunction) Arity() int {
	return n.arity
}

func (n *LoxNativeFunction) Call(interpreter *Interpreter, arguments []any) any {
	return n.fn(interpreter, arguments)
}

func (n *LoxNativeFunction) String() string {
	return n.name
}
