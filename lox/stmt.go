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

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

type Print struct {
	expression Expr
}

type Var struct {
	name        Token
	initializer Expr
}

type While struct {
	condition Expr
	body      Stmt
}

func (b Block) sealStmt()      {}
func (e Expression) sealStmt() {}
func (i If) sealStmt()         {}
func (p Print) sealStmt()      {}
func (v Var) sealStmt()        {}
func (w While) sealStmt()      {}
