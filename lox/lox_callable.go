package main

type LoxCallable struct {
	arity func() int
	call func(interpreter *Interpreter, arguments []any) any
	toString func() string
}

func (l *LoxCallable) String() string {
	return l.toString()
}
