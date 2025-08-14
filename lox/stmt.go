package main

type Stmt interface {
	sealStmt()
}

type Expression struct {
	expression Expr
}

type Print struct {
	expression Expr
}

type Var struct {
	name        Token
	initializer Expr
}

func (e Expression) sealStmt() {}
func (p Print) sealStmt()      {}
func (v Var) sealStmt()        {}
