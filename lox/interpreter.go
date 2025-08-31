package main

import (
	"fmt"
	"strings"
	"time"
)

type Interpreter struct {
	globals     *Environment
	environment *Environment
	// I would like to make the map keys *Expr, however, this seems to be disallowed
	// by Go. Even if I implement each concrete struct of Expr as pointer
	// receivers, that only makes e.g. *Assign be able to pass as Expr,
	// rather than a *Assign being able to pass as *Expr.
	locals map[Expr]int
}

type ReturnedValue struct {
	value any
}

var interpreter Interpreter = NewInterpreter()

func NewInterpreter() Interpreter {
	var environment = NewEnvironment()
	result := Interpreter{
		globals:     &environment,
		environment: &environment,
		locals:      make(map[Expr]int),
	}

	result.globals.define("clock", &LoxNativeFunction{
		arity: 0,
		fn: func(interpreter *Interpreter, arguments []any) any {
			return float64(time.Now().UnixMilli()) / 1000.0
		},
		name: "<native fn>",
	})

	return result
}

func (i *Interpreter) interpret_function_stmt(stmt Function) error {
	function := &LoxFunction{stmt, i.environment}
	i.environment.define(stmt.name.lexeme, function)
	return nil
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			if rte, ok := err.(RuntimeError); ok {
				runtimeError(rte)
			}
		}
	}
}

func (i *Interpreter) interpret_call_expr(expr *Call) (any, error) {
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

	if len(arguments) != function.Arity() {
		return nil, RuntimeError{expr.paren, fmt.Sprintf(
			"Expected %d arguments but got %d.",
			function.Arity(),
			len(arguments),
		)}
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) interpret_logical_expr(expr *Logical) (any, error) {
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

func (i *Interpreter) execute(stmt Stmt) (*ReturnedValue, error) {
	switch v := stmt.(type) {
	case If:
		return i.interpret_if_stmt(v)
	case Print:
		return nil, i.interpret_print_stmt(v)
	case Expression:
		return nil, i.interpret_expression_stmt(v)
	case Var:
		return nil, i.interpret_var_stmt(v)
	case Block:
		return i.interpret_block_stmt(v)
	case While:
		return i.interpret_while_stmt(v)
	case Function:
		return nil, i.interpret_function_stmt(v)
	case Return:
		return i.interpret_return_stmt(v) // This one actually returns a value
	case Class:
		return nil, i.interpret_class_stmt(v)
	default:
		panic(fmt.Sprintf("Unreachable. stmt has value %v; its type is %T which we don't know how to handle.", stmt, stmt))
	}
}

func (r *Resolver) resolve(expr Expr, depth int) {
	r.interpreter.locals[expr] = depth
}

func (i *Interpreter) interpret_while_stmt(stmt While) (*ReturnedValue, error) {
	for {
		cond, err := i.evaluate(stmt.condition)
		if err != nil {
			return nil, err
		}
		if !isTruthy(cond) {
			break
		}
		var res *ReturnedValue
		res, err = i.execute(stmt.body)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, nil
}

func (i *Interpreter) interpret_block_stmt(stmt Block) (*ReturnedValue, error) {
	innerEnv := NewEnvironment()
	innerEnv.enclosing = i.environment
	res, err := i.executeBlock(stmt.statements, &innerEnv)
	return res, err
}

func (i *Interpreter) interpret_class_stmt(stmt Class) error {
	i.environment.define(stmt.name.lexeme, nil)

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.methods {
		function := &LoxFunction{method, i.environment}
		methods[method.name.lexeme] = function
	}

	klass := &LoxClass{stmt.name.lexeme, methods}
	err := i.environment.assign(stmt.name, klass)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) interpret_if_stmt(stmt If) (*ReturnedValue, error) {
	cond, err := i.evaluate(stmt.condition)
	if err != nil {
		return nil, err
	}
	if isTruthy(cond) {
		res, err := i.execute(stmt.thenBranch)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	} else if stmt.elseBranch != nil { // TODO: this check might be bad, because I don't think elseBranch will be nil, it will be an empty Stmt
		res, err := i.execute(stmt.elseBranch)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, nil
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) (*ReturnedValue, error) {
	var result *ReturnedValue
	var err error
	previous := i.environment
	i.environment = environment
	defer func() {
		i.environment = previous
	}()
	for _, statement := range statements {
		result, err = i.execute(statement)
		if err != nil {
			return nil, err
		}
		if result != nil {
			return result, nil
		}
	}
	return nil, nil
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

func (i *Interpreter) interpret_assign_expr(expr *Assign) (any, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}

	distance, ok := i.locals[expr]
	if ok {
		i.environment.assignAt(distance, expr.name, value)
	} else {
		err = i.globals.assign(expr.name, value)
	}
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

	return RuntimeError{name, fmt.Sprintf("Inside assign: Undefined variable %q", name.lexeme)}
}

func (i *Interpreter) interpret_variable_expr(expr *Variable) (any, error) {
	return i.lookUpVariable(expr.name, expr)
}

func (i *Interpreter) lookUpVariable(name Token, expr Expr) (any, error) {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.getAt(distance, name.lexeme), nil
	} else {
		return i.globals.get(name)
	}
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	switch v := expr.(type) {
	case *Get:
		return i.interpret_get_expr(v)
	case *Set:
		return i.interpret_set_expr(v)
	case *This:
		return i.interpret_this_expr(v)
	case *Logical:
		return i.interpret_logical_expr(v)
	case *Binary:
		return i.interpret_binary_expr(v)
	case *Grouping:
		return i.interpret_grouping_expr(v)
	case *Literal:
		return interpret_literal_expr(v)
	case *Unary:
		return i.interpret_unary_expr(v)
	case *Variable:
		return i.interpret_variable_expr(v)
	case *Assign:
		return i.interpret_assign_expr(v)
	case *Call:
		return i.interpret_call_expr(v)
	default:
		panic(fmt.Sprintf("Unreachable. expr has value %v; its type is %T which we don't know how to handle.", expr, expr))
	}
}

func (i *Interpreter) interpret_set_expr(expr *Set) (any, error) {
	object, err := i.evaluate(expr.object)
	if err != nil {
		return nil, err
	}

	if inst, ok := object.(*LoxInstance); !ok {
		return nil, RuntimeError{expr.name, "Only instances have fields."}
	} else {
		value, err := i.evaluate(expr.value)
		if err != nil {
			return nil, err
		}
		inst.set(expr.name, value)
		return value, nil
	}
}

func (i *Interpreter) interpret_this_expr(expr *This) (any, error) {
	return i.lookUpVariable(expr.keyword, expr)
}

func (i *Interpreter) interpret_get_expr(expr *Get) (any, error) {
	object, err := i.evaluate(expr.object)
	if err != nil {
		return nil, err
	}
	if inst, ok := object.(*LoxInstance); ok {
		val, err := inst.get(expr.name)
		if err != nil {
			return nil, err
		}
		return val, nil
	}

	return nil, RuntimeError{expr.name, "Only instances have properties."}
}

func (i *Interpreter) interpret_binary_expr(expr *Binary) (any, error) {
	left, _ := i.evaluate(expr.left)
	right, _ := i.evaluate(expr.right)

	leftRV, leftIsRV := left.(*ReturnedValue)
	if leftIsRV {
		left = leftRV.value
	}
	rightRV, rightIsRV := right.(*ReturnedValue)
	if rightIsRV {
		right = rightRV.value
	}

	switch expr.operator.token_type {
	// TODO: fix the rest of these to return the error early if checkNumberOperands fails. If we return the err on the same line as the result, then we actually panic if left/right are not numbers, which is not what we want.
	case GREATER:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) > right.(float64), err
	case GREATER_EQUAL:
		err := checkNumberOperands(expr.operator, left, right)
		return left.(float64) >= right.(float64), err
	case LESS:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		err := checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
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

func (i *Interpreter) interpret_grouping_expr(expr *Grouping) (any, error) {
	result, err := i.evaluate(expr.expression)
	return result, err
}

func interpret_literal_expr(expr *Literal) (any, error) {
	return expr.value, nil
}

func (i *Interpreter) interpret_unary_expr(expr *Unary) (any, error) {
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

func (i *Interpreter) interpret_return_stmt(stmt Return) (*ReturnedValue, error) {
	var value any = nil
	var err error
	if stmt.value != nil {
		value, err = i.evaluate(stmt.value)
		if err != nil {
			return nil, err
		}
	}
	return &ReturnedValue{value}, nil
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}
	objectAsRV, objectIsRV := object.(*ReturnedValue)
	if objectIsRV {
		object = objectAsRV.value
	}
	switch v := object.(type) {
	case float64:
		text := fmt.Sprintf("%v", v)
		if strings.HasSuffix(text, ".0") {
			text = text[0 : len(text)-2]
		}
		return text
	default:
		if stringer, ok := (v).(fmt.Stringer); ok {
			return stringer.String()
		}
		return fmt.Sprintf("%v", v)
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
