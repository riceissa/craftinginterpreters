package main

type LoxClass struct {
	name string
}

func (l *LoxClass) String() string {
	return l.name
}

func (l *LoxClass) Arity() int {
	return 0
}

func (l *LoxClass) Call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(l)
	return instance, nil
}
