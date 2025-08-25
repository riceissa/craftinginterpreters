package main

type FunctionType int

const (
	FT_NONE = iota
	FT_FUNCTION
)

type Resolver struct {
	interpreter *Interpreter
	scopes []map[string]bool
	currentFunction FunctionType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes: nil,
		currentFunction: FT_NONE,
	}
}

func (r *Resolver) resolveBlockStmt(stmt Block) {
	beginScope()
	resolveStatements(stmt.statements)
	endScope()
}

func resolveFunctionStmt(stmt Function) {
	declare(stmt.name)
	define(stmt.name)

	resolveFunction(stmt, FT_FUNCTION)
}

func (r *Resolver) resolveVarStmt(stmt Var) {
	declare(stmt.name)
	if stmt.initializer != nil {
		resolve(stmt.initializer)
	}
	define(stmt.name)
}

func resolveAssignExpr(expr Assign) {
	resolve(expr.value)
	resolveLocal(expr, expr.name)
}

func (r *Resolver) resolveVariableExpr(expr Variable) {
	if len(scopes) > 0 && !scopes[len(scopes) - 1][expr.name.lexeme] {
		log_parse_error(expr.name, "Can't read local variable in its own initializer: %q\n")
	}

	resolveLocal(expr, expr.name)
}

func resolveStatements(statements []Stmt) {
	for statement := range statements {
		resolveStatement(statement)
	}
}

func resolveStatement(stmt Stmt) {
}

func resolveExpr(expr Expr) {
}

func beginScope() {
	scopes = append(scopes, make(map[string]bool))
}

func endScope() {
	scopes = scopes[:len(scopes)-1]
}

func declare(name Token) {
	if len(scopes) == 0 {
		return
	}

	scope := scopes[len(scopes) - 1]
	if _, ok := scope[name.lexeme]; ok {
		log_parse_error(name, "Already a variable with this name in this scope: %q\n")
	}
	scope[name.lexeme] = false
}

func define(name Token) {
	if len(scopes) == 0 {
		return
	}
	scopes[len(scopes) - 1][name.lexeme] = true
}

func resolveLocal(expr Expr, name Token) {
	for i := len(scopes) - 1; i >= 0; i-- {
		if _, ok := scopes[i][name.lexeme]; ok {
			interpreter.resolve(expr, len(scopes) - 1 - i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function Function, ft FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = ft

	beginScope()
	for param := range function.params {
		declare(param)
		define(param)
	}
	resolve(function.body)
	endScope()
	r.currentFunction = enclosingFunction
}

func resolveExpressionStmt(stmt Expression) {
	resolve(stmt.expression)
}

func resolveIfStmt(stmt If) {
	resolve(stmt.condition)
	resolve(stmt.thenBranch)
	if stmt.elseBranch != nil {
		resolve(stmt.elseBranch)
	}
}

func resolvePrintStmt(stmt Print) {
	resolve(stmt.expression)
}

func resolveReturnStmt(stmt Return) {
	if r.currentFunction == FT_NONE {
		log_parse_error(stmt.keyword, "Can't return from top-level code: %q\n")
	}

	if stmt.value != nil {
		resolve(stmt.value)
	}
}

func resolveWhileStmt(stmt While) {
	resolve(stmt.condition)
	resolve(stmt.body)
}

func resolveBinaryExpr(expr Binary) {
	resolve(expr.left)
	resolve(expr.right)
}

func resolveCallExpr(expr Call) {
	resolve(expr.callee)

	for argument := range expr.arguments {
		resolve(argument)
	}
}

func resolveGroupingExpr(expr Grouping) {
	resolve(expr.expression)
}

func resolveLiteralExpr(expr Literal) {
	// nothing to do
}

func resolveLogicalExpr(expr Logical) {
	resolve(expr.left)
	resolve(expr.right)
}

func resolveUnaryExpr(expr Unary) {
	resolve(expr.right)
}
