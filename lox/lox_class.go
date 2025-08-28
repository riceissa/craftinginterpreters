package main

type LoxClass struct {
	name string
}

func (l LoxClass) String() string {
	return l.name
}
