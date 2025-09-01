package main

type LoxClass struct {
	name    string
	superclass *LoxClass
	methods map[string]*LoxFunction
}

func (l *LoxClass) findMethod(name string) *LoxFunction {
	if v, ok := l.methods[name]; ok {
		return v
	}
	if l.superclass != nil {
		return l.superclass.findMethod(name)
	}
	return nil
}

func (l *LoxClass) String() string {
	return l.name
}

func (l *LoxClass) Arity() int {
	initializer := l.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (l *LoxClass) Call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(l)
	initializer := l.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).Call(interpreter, arguments)
	}

	return instance, nil
}
