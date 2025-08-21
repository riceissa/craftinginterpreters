package main

import (
	"fmt"
	"strings"
	"time"
)

type Interpreter struct {
	globals     *Environment
	environment *Environment
}

var interpreter Interpreter = NewInterpreter()

func NewInterpreter() Interpreter {
	var environment = NewEnvironment()
	result := Interpreter{
		globals:     &environment,
		environment: &environment,
	}

	result.globals.define("clock", LoxCallable{
		arity: func() int { return 0 },
		call: func(interpreter *Interpreter, arguments []any) any {
			return float64(time.Now().UnixMilli()) / 1000.0
		},
		toString: func() string { return "<native fn>" },
	})

	return result
}

func (e *Environment) interpret_function_stmt(stmt Stmt) {
	fun, ok := stmt.(Function)
	if !ok {
		panic("was expecting a Function here")
	}
	function := LoxFunction{fun}
	e.define(fun.name.lexeme, function)
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			if rte, ok := err.(RuntimeError); ok {
				runtimeError(rte)
			}
		}
	}
}

func (i *Interpreter) interpret_call_expr(expr Call) (any, error) {
	callee, err := i.evaluate(expr.callee)
	if err != nil {
		return nil, err
	}
	arguments := []any{}
	for _, argument := range expr.arguments {
		res, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, res)
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, RuntimeError{expr.paren, "Can only call functions and classes."}
	}

	if len(arguments) != function.arity() {
		return nil, RuntimeError{expr.paren, fmt.Sprintf(
			"Expected %d arguments but got %d.",
			function.arity(),
			len(arguments),
		)}
	}

	return function.call(i, arguments), nil
}

func (i *Interpreter) interpret_logical_expr(expr Logical) (any, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}

	if expr.operator.token_type == OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.right)
}

func (i *Interpreter) execute(stmt Stmt) error {
	switch v := stmt.(type) {
	case If:
		return i.interpret_if_stmt(v)
	case Print:
		return i.interpret_print_stmt(v)
	case Expression:
		return i.interpret_expression_stmt(v)
	case Var:
		return i.interpret_var_stmt(v)
	case Block:
		return i.interpret_block_stmt(v)
	case While:
		return i.interpret_while_stmt(v)
	default:
		panic(fmt.Sprintf("Unreachable. stmt has value %v; its type is %T which we don't know how to handle.", stmt, stmt))
	}
}

func (i *Interpreter) interpret_while_stmt(stmt While) error {
	for {
		cond, err := i.evaluate(stmt.condition)
		if err != nil {
			return err
		}
		if !isTruthy(cond) {
			break
		}
		err = i.execute(stmt.body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) interpret_block_stmt(stmt Block) error {
	innerEnv := NewEnvironment()
	innerEnv.enclosing = i.environment
	previous := i.environment
	i.environment = &innerEnv
	err := i.executeBlock(stmt.statements)
	i.environment = previous
	return err
}

func (i *Interpreter) interpret_if_stmt(stmt If) error {
	cond, err := i.evaluate(stmt.condition)
	if err != nil {
		return err
	}
	if isTruthy(cond) {
		err := i.execute(stmt.thenBranch)
		if err != nil {
			return err
		}
	} else if stmt.elseBranch != nil { // TODO: this check might be bad, because I don't think elseBranch will be nil, it will be an empty Stmt
		err := i.execute(stmt.elseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) executeBlock(statements []Stmt) error {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) interpret_var_stmt(stmt Var) error {
	var value any
	var err error
	if stmt.initializer != nil {
		value, err = i.evaluate(stmt.initializer)
		if err != nil {
			return err
		}
	}

	e := i.environment
	e.define(stmt.name.lexeme, value)
	return nil
}

func (i *Interpreter) interpret_assign_expr(expr Assign) (any, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	err = i.environment.assign(expr.name, value)
	return value, err
}

func (e *Environment) assign(name Token, value any) error {
	_, ok := e.values[name.lexeme]
	if ok {
		e.values[name.lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return nil
	}

	return RuntimeError{name, fmt.Sprintf("Undefined variable %q; %q", name.lexeme)}
}

func (i *Interpreter) interpret_variable_expr(expr Variable) (any, error) {
	return i.environment.get(expr.name)
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	switch v := expr.(type) {
	case Logical:
		return i.interpret_logical_expr(v)
	case Binary:
		return i.interpret_binary_expr(v)
	case Grouping:
		return i.interpret_grouping_expr(v)
	case Literal:
		return interpret_literal_expr(v)
	case Unary:
		return i.interpret_unary_expr(v)
	case Variable:
		return i.interpret_variable_expr(v)
	case Assign:
		return i.interpret_assign_expr(v)
	default:
		panic(fmt.Sprintf("Unreachable. expr has value %v; its type is %T which we don't know how to handle.", expr, expr))
	}
}

func (i *Interpreter) interpret_binary_expr(expr Binary) (any, error) {
	left, _ := i.evaluate(expr.left)
	right, _ := i.evaluate(expr.right)

	switch expr.operator.token_type {
	case GREATER:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) > right.(float64), err
	case GREATER_EQUAL:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) >= right.(float64), err
	case LESS:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) < right.(float64), err
	case LESS_EQUAL:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) <= right.(float64), err
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	case MINUS:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) - right.(float64), err
	case PLUS:
		leftFloat, leftIsFloat := left.(float64)
		rightFloat, rightIsFloat := right.(float64)
		if leftIsFloat && rightIsFloat {
			return leftFloat + rightFloat, nil
		}

		leftString, leftIsString := left.(string)
		rightString, rightIsString := right.(string)
		if leftIsString && rightIsString {
			return leftString + rightString, nil
		}

		return nil, RuntimeError{expr.operator, "Operands must be two numbers or two strings."}
	case SLASH:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) / right.(float64), err
	case STAR:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) * right.(float64), err
	}

	panic("Unreachable")
}

func (i *Interpreter) interpret_grouping_expr(expr Grouping) (any, error) {
	result, err := i.evaluate(expr.expression)
	return result, err
}

func interpret_literal_expr(expr Literal) (any, error) {
	return expr.value, nil
}

func (i *Interpreter) interpret_unary_expr(expr Unary) (any, error) {
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.token_type {
	case BANG:
		return !isTruthy(right), nil
	case MINUS:
		err := checkNumberOperand(expr.operator, right)
		return -(right.(float64)), err
	}

	panic("Unreachable")
}

func (i *Interpreter) interpret_expression_stmt(stmt Expression) error {
	_, err := i.evaluate(stmt.expression)
	return err
}

func (i *Interpreter) interpret_print_stmt(stmt Print) error {
	value, err := i.evaluate(stmt.expression)
	if err != nil {
		return err
	}
	fmt.Println(stringify(value))
	return nil
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}
	switch object.(type) {
	case float64:
		text := fmt.Sprintf("%v", object)
		if strings.HasSuffix(text, ".0") {
			text = text[0 : len(text)-2]
		}
		return text
	default:
		return fmt.Sprintf("%v", object)
	}
}

func checkNumberOperand(operator Token, operand any) error {
	switch operand.(type) {
	case float64:
		return nil
	default:
		return RuntimeError{operator, "Operand must be a number."}
	}
}

func checkNumberOperands(operator Token, left any, right any) error {
	_, leftIsFloat := left.(float64)
	_, rightIsFloat := right.(float64)
	if leftIsFloat && rightIsFloat {
		return nil
	}
	return RuntimeError{operator, "Operands must be numbers."}
}

func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	switch object.(type) {
	case bool:
		return object.(bool)
	default:
		return true
	}
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	// TODO: might need to use something like reflect.DeepEqual
	// or whatever; I don't know if Go's native == operator
	// is what we want here.
	return a == b
}
