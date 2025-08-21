package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration Function
}

func (f *LoxFunction) arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %v >", f.declaration.name.lexeme)
}

func (f *LoxFunction) call(interpreter Interpreter, arguments []any) any {
	environment := NewEnvironment()
	environment.enclosing = interpreter.globals // TODO: this might not be right. what's the parameter passed to Environment() constructor in the java version?
	for i := range len(f.declaration.params) {
		environment.define(f.declaration.params[i].lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaration.body, &environment)
	return nil
}
