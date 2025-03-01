package core

import (
	"fmt"
	"os"
)

type Function struct {
	fmt.Stringer
	Arity int
	Call  func(args []any) (any, error)
}

func (f Function) String() string {
	return fmt.Sprintf("f(%v)", f.Arity)
}

type Expr interface {
	fmt.Stringer
	Interpret(environment Environment) any
}

type BinaryExpr struct {
	Expr
	rightExpr Expr
	operator  Token
	leftExpr  Expr
}

type UnaryExpr struct {
	Expr
	operator Token
	expr     Expr
}

type GroupingExpr struct {
	Expr
	expr Expr
}

type LiteralExpr struct {
	Expr
	value Token
}

type AssignExpr struct {
	Expr
	name     Token
	expr     Expr
	operator Token
}

type BlockExpr struct {
	Expr
	program []Expr
}

type IfExpr struct {
	Expr
	condition  Expr
	thenBranch Expr
	elseBranch Expr
}

type WhileExpr struct {
	Expr
	condition  Expr
	loopBranch Expr
}

type InvalidExpr struct {
	Expr
}

type CallExpr struct {
	Expr
	f    Expr
	args []Expr
}

func (expr AssignExpr) String() string {
	op := ""
	if expr.operator.Type == COLON_EQ {
		op = ":="
	} else if expr.operator.Type == EQUAL {
		op = "="
	}
	return fmt.Sprintf("(%v %v %v)", expr.name.String(), op, expr.expr.String())
}
func (expr BinaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", expr.leftExpr.String(), expr.operator.String(), expr.rightExpr.String())
}
func (expr UnaryExpr) String() string {
	return fmt.Sprintf("(%v%v)", expr.operator.String(), expr.expr.String())
}
func (expr GroupingExpr) String() string {
	return fmt.Sprintf("(%v)", expr.expr.String())
}
func (expr LiteralExpr) String() string {
	return expr.value.String()
}
func (expr InvalidExpr) String() string {
	return "Invalid Expression"
}
func (expr BlockExpr) String() string {
	res := ""
	for _, expr := range expr.program {
		res += "\t" + expr.String()
	}
	return res
}
func (expr IfExpr) String() string {
	res := fmt.Sprintf("if (%v) {\n", expr.condition.String())
	res += expr.thenBranch.String()
	if expr.elseBranch != nil {
		res += fmt.Sprintf("} else {\n")
		res += expr.elseBranch.String()
	}
	res += "}"
	return res
}
func (expr WhileExpr) String() string {
	res := fmt.Sprintf("while (%v) {\n", expr.condition.String())
	res += expr.loopBranch.String()
	res += "}"
	return res
}
func (expr CallExpr) String() string {
	res := fmt.Sprintf("%s(", expr.f.String())
	for i, arg := range expr.args {
		if i > 0 {
			res += ", "
		}
		res += arg.String()
	}
	res += ")"
	return res
}

func (expr AssignExpr) Interpret(environment Environment) any {
	data := expr.expr.Interpret(environment)
	if expr.operator.Type == COLON_EQ {
		environment.declareVar(expr.name.String(), data)
	} else if expr.operator.Type == EQUAL {
		err := environment.setVar(expr.name.String(), data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		err := fmt.Errorf("Invalid assignment operator: %s", expr.operator.Type)
		fmt.Println(err)
		os.Exit(1)
	}
	return data
}

func (expr BinaryExpr) Interpret(environment Environment) any {
	l := expr.leftExpr.Interpret(environment)
	r := expr.rightExpr.Interpret(environment)

	switch expr.operator.Type {
	case PLUS:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt + rightInt
			}
		}
		return l.(float64) + r.(float64)
	case MINUS:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt - rightInt
			}
		}
		return l.(float64) - r.(float64)
	case STAR:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt * rightInt
			}
		}
		return l.(float64) * r.(float64)
	case SLASH:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt / rightInt
			}
		}
		return l.(float64) / r.(float64)
	case PERCENT:
		return l.(int64) % r.(int64)
	case EQUAL_EQ:
		return l == r
	case BANG_EQ:
		return l != r
	case GREATER:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt > rightInt
			}
		}
		return l.(float64) > r.(float64)
	case GREATER_EQ:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt >= rightInt
			}
		}
		return l.(float64) >= r.(float64)
	case LESS:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt < rightInt
			}
		}
		return l.(float64) < r.(float64)
	case LESS_EQ:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt <= rightInt
			}
		}
		return l.(float64) <= r.(float64)
	case BAR:
		return l.(int64) | r.(int64)
	case BAR_BAR:
		return l.(bool) || r.(bool)
	case AMP:
		return l.(int64) & r.(int64)
	case AMP_AMP:
		return l.(bool) && r.(bool)
	default:
		return nil
	}
}

func (expr UnaryExpr) Interpret(environment Environment) any {
	res := expr.expr.Interpret(environment)
	switch expr.operator.Type {
	case BANG:
		return !res.(bool)
	case MINUS:
		if _, ok := res.(int64); ok {
			return -res.(int64)
		}
		return -res.(float64)
	default:
		return nil
	}
}

func (expr GroupingExpr) Interpret(environment Environment) any {
	return expr.expr.Interpret(environment)
}

func (expr LiteralExpr) Interpret(environment Environment) any {
	switch expr.value.Type {
	case TRUE:
		return true
	case FALSE:
		return false
	case IDENTIFIER:
		v, err := environment.findVar(expr.value.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return nil
		}
		return v
	default:
		return expr.value.Value
	}
}

func (expr InvalidExpr) Interpret(environment Environment) any {
	return nil
}

func (expr BlockExpr) Interpret(environment Environment) any {
	var res any
	environment.push(make(map[string]any))
	for _, expr := range expr.program {
		res = expr.Interpret(environment)
	}
	environment.pop()
	return res
}

func (expr IfExpr) Interpret(environment Environment) any {
	if expr.condition.Interpret(environment).(bool) {
		return expr.thenBranch.Interpret(environment)
	} else if expr.elseBranch != nil {
		return expr.elseBranch.Interpret(environment)
	}

	return nil
}

func (expr WhileExpr) Interpret(environment Environment) any {
	for expr.condition.Interpret(environment).(bool) {
		expr.loopBranch.Interpret(environment)
	}

	return nil
}

func (expr CallExpr) Interpret(environment Environment) any {
	f, err := environment.findVar(expr.f.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil
	}

	function, ok := f.(Function)
	if !ok {
		fmt.Println("Expected function")
		os.Exit(1)
		return nil
	}
	args := []any{}

	for _, arg := range expr.args {
		args = append(args, arg.Interpret(environment))
	}
	value, err := function.Call(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil
	}
	return value
}
