package main

type Stmt interface {
	sealStmt()
}

type Block struct {
	statements []Stmt
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

func (b Block) sealStmt() {}
func (e Expression) sealStmt() {}
func (p Print) sealStmt()      {}
func (v Var) sealStmt()        {}
