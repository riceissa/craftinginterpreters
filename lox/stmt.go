package main

type Stmt interface {
	stmtSealer()
}

type Expression struct {
	expression Expr
}

type Print struct {
	expression Expr
}

func (e Expression) stmtSealer() {}
func (p Print) stmtSealer() {}
