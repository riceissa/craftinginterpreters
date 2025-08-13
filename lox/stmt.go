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

type Var struct {
	name Token
	initializer Expr
}

func (e Expression) stmtSealer() {}
func (p Print) stmtSealer() {}
func (v Var) stmtSealer() {}
