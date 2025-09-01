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
	sealExpr()
}

type Assign struct {
	name  Token
	value Expr
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

type Get struct {
	object Expr
	name   Token
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value any
}

type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

type Set struct {
	object Expr
	name   Token
	value  Expr
}

type Super struct {
	keyword Token
	method  Token
}

type This struct {
	keyword Token
}

type Unary struct {
	operator Token
	right    Expr
}

type Variable struct {
	name Token
}

func (b *Binary) sealExpr()   {}
func (c *Call) sealExpr()     {}
func (g *Get) sealExpr()      {}
func (g *Grouping) sealExpr() {}
func (l *Literal) sealExpr()  {}
func (l *Logical) sealExpr()  {}
func (s *Set) sealExpr()      {}
func (s *Super) sealExpr()    {}
func (t *This) sealExpr()     {}
func (u *Unary) sealExpr()    {}
func (v *Variable) sealExpr() {}
func (a *Assign) sealExpr()   {}

// Assert we've correctly implemented the interface.
var _ Expr = &Binary{}
var _ Expr = &Call{}
var _ Expr = &Get{}
var _ Expr = &Grouping{}
var _ Expr = &Literal{}
var _ Expr = &Logical{}
var _ Expr = &Set{}
var _ Expr = &Super{}
var _ Expr = &This{}
var _ Expr = &Unary{}
var _ Expr = &Variable{}
var _ Expr = &Assign{}
