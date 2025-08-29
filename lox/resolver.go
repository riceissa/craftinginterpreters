package main

type FunctionType int

const (
	FT_NONE = iota
	FT_FUNCTION
	FT_METHOD
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          nil,
		currentFunction: FT_NONE,
	}
}

func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	switch v := stmt.(type) {
	case Class:
		r.resolveClassStmt(v)
	case Block:
		r.resolveBlockStmt(v)
	case Function:
		r.resolveFunctionStmt(v)
	case Var:
		r.resolveVarStmt(v)
	case Expression:
		r.resolveExpressionStmt(v)
	case If:
		r.resolveIfStmt(v)
	case Print:
		r.resolvePrintStmt(v)
	case Return:
		r.resolveReturnStmt(v)
	case While:
		r.resolveWhileStmt(v)
	default:
		panic("Unreachable.")
	}
}

func (r *Resolver) resolveExpr(expr Expr) {
	switch v := expr.(type) {
	case *Get:
		r.resolveGetExpr(v)
	case *Set:
		r.resolveSetExpr(v)
	case *Assign:
		r.resolveAssignExpr(v)
	case *Variable:
		r.resolveVariableExpr(v)
	case *Binary:
		r.resolveBinaryExpr(v)
	case *Call:
		r.resolveCallExpr(v)
	case *Grouping:
		r.resolveGroupingExpr(v)
	case *Literal:
		r.resolveLiteralExpr(v)
	case *Logical:
		r.resolveLogicalExpr(v)
	case *Unary:
		r.resolveUnaryExpr(v)
	default:
		panic("Unreachable.")
	}
}

func (r *Resolver) resolveSetExpr(expr *Set) {
	r.resolveExpr(expr.value)
	r.resolveExpr(expr.object)
}

func (r *Resolver) resolveGetExpr(expr *Get) {
	r.resolveExpr(expr.object)
}

func (r *Resolver) resolveBlockStmt(stmt Block) {
	r.beginScope()
	r.resolveStatements(stmt.statements)
	r.endScope()
}

func (r *Resolver) resolveClassStmt(stmt Class) {
	r.declare(stmt.name)
	r.define(stmt.name)

	for _, method := range stmt.methods {
		var declaration FunctionType = FT_METHOD
		r.resolveFunction(method, declaration)
	}
}

func (r *Resolver) resolveFunctionStmt(stmt Function) {
	r.declare(stmt.name)
	r.define(stmt.name)

	r.resolveFunction(stmt, FT_FUNCTION)
}

func (r *Resolver) resolveVarStmt(stmt Var) {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolveExpr(stmt.initializer)
	}
	r.define(stmt.name)
}

func (r *Resolver) resolveAssignExpr(expr *Assign) {
	r.resolveExpr(expr.value)
	r.resolveLocal(expr, expr.name)
}

func (r *Resolver) resolveVariableExpr(expr *Variable) {
	if len(r.scopes) > 0 {
		v, ok := r.scopes[len(r.scopes)-1][expr.name.lexeme]
		if ok && !v {
			log_parse_error(expr.name, "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.name)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.lexeme]; ok {
		log_parse_error(name, "Already a variable with this name in this scope")
	}
	scope[name.lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.lexeme]; ok {
			r.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function Function, ft FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = ft

	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) resolveExpressionStmt(stmt Expression) {
	r.resolveExpr(stmt.expression)
}

func (r *Resolver) resolveIfStmt(stmt If) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolveStmt(stmt.elseBranch)
	}
}

func (r *Resolver) resolvePrintStmt(stmt Print) {
	r.resolveExpr(stmt.expression)
}

func (r *Resolver) resolveReturnStmt(stmt Return) {
	if r.currentFunction == FT_NONE {
		log_parse_error(stmt.keyword, "Can't return from top-level code")
	}

	if stmt.value != nil {
		r.resolveExpr(stmt.value)
	}
}

func (r *Resolver) resolveWhileStmt(stmt While) {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
}

func (r *Resolver) resolveBinaryExpr(expr *Binary) {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
}

func (r *Resolver) resolveCallExpr(expr *Call) {
	r.resolveExpr(expr.callee)

	for _, argument := range expr.arguments {
		r.resolveExpr(argument)
	}
}

func (r *Resolver) resolveGroupingExpr(expr *Grouping) {
	r.resolveExpr(expr.expression)
}

func (r *Resolver) resolveLiteralExpr(expr *Literal) {
	// nothing to do
}

func (r *Resolver) resolveLogicalExpr(expr *Logical) {
	r.resolveExpr(expr.left)
	r.resolveExpr(expr.right)
}

func (r *Resolver) resolveUnaryExpr(expr *Unary) {
	r.resolveExpr(expr.right)
}
