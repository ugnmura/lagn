package core

import (
	"fmt"
	"os"
)

type Environment []map[string]any

func (env *Environment) push(v map[string]any) Environment {
	*env = append(*env, v)
	return *env
}

func (env *Environment) pop() map[string]any {
	l := len(*env)
	if l == 0 {
		panic("Cannot pop from an empty environment")
	}

	res := (*env)[l-1]
	*env = (*env)[:l-1]
	return res
}

func (env Environment) findVar(name string) (any, error) {
	for i := len(env) - 1; i >= 0; i-- {
		if val, ok := env[i][name]; ok {
			return val, nil
		}
	}

	return nil, fmt.Errorf("Variable %s not found", name)
}

func (env *Environment) setVar(name string, value any) error {
	for i := len(*env) - 1; i >= 0; i-- {
		if _, ok := (*env)[i][name]; ok {
			(*env)[i][name] = value
			return nil
		}
	}

	return fmt.Errorf("Variable %s not found", name)
}

func (env *Environment) declareVar(name string, value any) {
	l := len(*env)
	(*env)[l-1][name] = value
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
	left := expr.leftExpr.Interpret(environment)
	right := expr.rightExpr.Interpret(environment)
	l, _ := left.(float64)
	r, _ := right.(float64)

	switch expr.operator.Type {
	case PLUS:
		return l + r
	case MINUS:
		return l - r
	case STAR:
		return l * r
	case SLASH:
		return l / r
	case EQUAL_EQ:
		return left == right
	case BANG_EQ:
		return left != right
	case GREATER:
		return l > r
	case GREATER_EQ:
		return l >= r
	case LESS:
		return l < r
	case LESS_EQ:
		return l <= r
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
