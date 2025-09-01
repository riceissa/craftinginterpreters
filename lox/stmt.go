package main

type Stmt interface {
	sealStmt()
}

type Block struct {
	statements []Stmt
}

type Class struct {
	name       Token
	superclass *Variable
	methods    []Function
}

type Expression struct {
	expression Expr
}

type Function struct {
	name   Token
	params []Token
	body   []Stmt
}

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

type Print struct {
	expression Expr
}

type Return struct {
	keyword Token
	value   Expr
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
func (c Class) sealStmt()      {}
func (e Expression) sealStmt() {}
func (f Function) sealStmt()   {}
func (i If) sealStmt()         {}
func (p Print) sealStmt()      {}
func (r Return) sealStmt()     {}
func (v Var) sealStmt()        {}
func (w While) sealStmt()      {}

// Assert that we've implemented the interface
var _ Stmt = Block{}
var _ Stmt = Class{}
var _ Stmt = Expression{}
var _ Stmt = Function{}
var _ Stmt = If{}
var _ Stmt = Print{}
var _ Stmt = Return{}
var _ Stmt = Var{}
var _ Stmt = While{}
