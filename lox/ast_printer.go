package main

func print(expr Expr) string {
	switch v := expr.(type) {
	case Binary:
		parenthesize(v.operator.lexeme, v.left, v.right)
	case Grouping:
		parenthesize("group", v.expression)
	case Literal:
		if v.value == nil {
			return "nil"
		}
		return string(v.value)
	case Unary:
		parenthesize(v.operator.lexeme, v.right)
	default:
	}
}

func parenthesize(name string, exprs ...Expr) {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr)
	}
	builder.WriteString(")")
	return builder.String()
}
