package main

type Expr interface {
	// Go doesn't have unions/sum types. So we create an interface
	// (Expr) that our various structs will all "implement".  But in
	// order to do this, we need this garbage method that is used to
	// "seal" the interface so that random structs won't be
	// considered to be "implementing" this interface.  (If the Expr
	// interface had no methods, then it would be the empty
	// interface, so then automatically every struct would be
	// considered to be "implementing" the interface!  We don't want
	// that.)
	exprSealer()
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value any
}

type Unary struct {
	operator Token
	right    Expr
}

type Variable struct {
	name Token
}

func (b Binary) exprSealer()   {}
func (g Grouping) exprSealer() {}
func (l Literal) exprSealer()  {}
func (u Unary) exprSealer()    {}
func (v Variable) exprSealer() {}
