package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration Function
	closure     *Environment
	isInitializer bool
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

	if f.isInitializer {
		return f.closure.getAt(0, "this")
	}

	return nil, nil
}

func (l *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment()
	environment.enclosing = l.closure
	environment.define("this", instance)
	return &LoxFunction{l.declaration, &environment, l.isInitializer}
}

type LoxNativeFunction struct {
	arity int
	fn    func(*Interpreter, []any) any
	name  string
}

func (n *LoxNativeFunction) Arity() int {
	return n.arity
}

func (n *LoxNativeFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return n.fn(interpreter, arguments), nil
}

func (n *LoxNativeFunction) String() string {
	return n.name
}

// Assert that LoxFunction and LoxNativeFunction implement the LoxCallable
// interface; perhaps this should go in separate testing code so it doesn't
// happen every time we compile the code, but since the point of this project
// isn't to be the most efficient interpreter, but to just learn how to write
// one, I am fine with doing this. It will probably save some painful debugging
// sessions in the future, although I have gotten reasonably good at intuiting
// when the source of a bug is that the thing I thought was implementing the
// interface no longer is (or never was).
var _ LoxCallable = &LoxFunction{}
var _ LoxCallable = &LoxNativeFunction{}
