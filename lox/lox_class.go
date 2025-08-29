package main

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func (l *LoxClass) findMethod(name string) *LoxFunction {
	return l.methods[name]
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
