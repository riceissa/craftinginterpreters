package main

type LoxInstance struct {
	klass *LoxClass
}

func (l *LoxInstance) String() string {
	return l.klass.name + " instance"
}
