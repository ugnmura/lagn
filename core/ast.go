package core

import (
	"fmt"
)

type Function struct {
	fmt.Stringer
	Arity int
	Call  func(env Environment, args []any) (any, error)
}

func (f Function) String() string {
	return fmt.Sprintf("f(%v)", f.Arity)
}

type Expr interface {
	fmt.Stringer
	Interpret(environment Environment) (any, error)
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


type CallExpr struct {
	Expr
	f    Expr
	args []Expr
}

type FnDeclExpr struct {
  Expr
  name Token
	args []Token
  program Expr
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
func (expr FnDeclExpr) String() string {
	res := fmt.Sprintf("%s = (", expr.name.String())
	for i, arg := range expr.args {
		if i > 0 {
			res += ", "
		}
		res += arg.String()
	}
	res += ") => \n"
  res += expr.program.String()
	return res
}

func (expr AssignExpr) Interpret(environment Environment) (any, error) {
	data, err := expr.expr.Interpret(environment)
  if err != nil {
    return nil, err
  }
	if expr.operator.Type == COLON_EQ {
		environment.declareVar(expr.name.String(), data)
	} else if expr.operator.Type == EQUAL {
		err := environment.setVar(expr.name.String(), data)
		if err != nil {
      return nil, fmt.Errorf("Invalid assignment operator: %s", expr.operator.Type)
		}
	} else {
		return nil, fmt.Errorf("Invalid assignment operator: %s", expr.operator.Type)
	}
	return data, nil
}

func (expr BinaryExpr) Interpret(environment Environment) (any, error) {
	l, err := expr.leftExpr.Interpret(environment)
  if err != nil {
    return nil, err
  }
	r, err := expr.rightExpr.Interpret(environment)
  if err != nil {
    return nil, err
  }

	switch expr.operator.Type {
  case PLUS:
      switch left := l.(type) {
      case int64:
          switch right := r.(type) {
          case int64:
              return left + right, nil
          case float64:
              return float64(left) + right, nil
          case string:
              return fmt.Sprintf("%d%s", left, right), nil
          default:
              return nil, fmt.Errorf("Unsupported type for addition: Int + %T", right)
          }
      case float64:
          switch right := r.(type) {
          case int64:
              return left + float64(right), nil
          case float64:
              return left + right, nil
          case string:
              return fmt.Sprintf("%f%s", left, right), nil
          default:
              return nil, fmt.Errorf("Unsupported type for addition: Float + %T", right)
          }
      case string:
          switch right := r.(type) {
          case int64:
              return left + fmt.Sprintf("%d", right), nil
          case float64:
              return left + fmt.Sprintf("%f", right), nil
          case string:
              return left + right, nil
          default:
              return nil, fmt.Errorf("Unsupported type for addition: String + %T", right)
          }
      default:
          return nil, fmt.Errorf("Unsupported type for addition: %T + %T", l, r)
      }
	case MINUS:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt - rightInt, nil
			}
		}
		return l.(float64) - r.(float64), nil
	case STAR:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt * rightInt, nil
			}
		}
		return l.(float64) * r.(float64), nil
	case SLASH:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt / rightInt, nil
			}
		}
		return l.(float64) / r.(float64), nil
	case PERCENT:
		return l.(int64) % r.(int64), nil
	case EQUAL_EQ:
		return l == r, nil
	case BANG_EQ:
		return l != r, nil
	case GREATER:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt > rightInt, nil
			}
		}
		return l.(float64) > r.(float64), nil
	case GREATER_EQ:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt >= rightInt, nil
			}
		}
		return l.(float64) >= r.(float64), nil
	case LESS:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt < rightInt, nil
			}
		}
		return l.(float64) < r.(float64), nil
	case LESS_EQ:
		if leftInt, ok := l.(int64); ok {
			if rightInt, ok := r.(int64); ok {
				return leftInt <= rightInt, nil
			}
		}
		return l.(float64) <= r.(float64), nil
	case BAR:
		return l.(int64) | r.(int64), nil
	case BAR_BAR:
		return l.(bool) || r.(bool), nil
	case AMP:
		return l.(int64) & r.(int64), nil
	case AMP_AMP:
		return l.(bool) && r.(bool), nil
	default:
		return nil, fmt.Errorf("Invalid Binary Operator %v", expr.operator)
	}
}

func (expr UnaryExpr) Interpret(environment Environment) (any, error) {
	res, err := expr.expr.Interpret(environment)
  if err != nil {
    return nil, err
  }
	switch expr.operator.Type {
	case BANG:
		if r, ok := res.(bool); ok {
      return !r, nil
    } 
    return nil, fmt.Errorf("Expected Type bool, got %T", res)
	case MINUS:
		if r, ok := res.(int64); ok {
			return -r, nil
		}
		if r, ok := res.(float64); ok {
			return -r, nil
		}
    return nil, fmt.Errorf("Expected number, got %T", res)
	default:
		return nil, fmt.Errorf("Invalid Unary Operator %v", expr.operator.Type)
	}
}

func (expr GroupingExpr) Interpret(environment Environment) (any, error) {
	return expr.expr.Interpret(environment)
}

func (expr LiteralExpr) Interpret(environment Environment) (any, error) {
	switch expr.value.Type {
	case TRUE:
		return true, nil
	case FALSE:
		return false, nil
	case IDENTIFIER:
		v, err := environment.findVar(expr.value.String())
		if err != nil {
      return nil, err
		}
		return v, nil
	default:
		return expr.value.Value, nil
	}
}

func (expr BlockExpr) Interpret(environment Environment) (any, error) {
	var res any
  var err error
	environment.push(make(map[string]any))
	for _, expr := range expr.program {
		res, err = expr.Interpret(environment)

    if err != nil {
      return nil, err
    }
	}
	environment.pop()
	return res, nil
}

func (expr IfExpr) Interpret(environment Environment) (any, error) {
  condVal, err := expr.condition.Interpret(environment)
  if err != nil {
    return nil, err 
  }
  if _, ok := condVal.(bool); !ok {
    return nil, fmt.Errorf("Expected Type bool, got %T", condVal)
  }

	if condVal.(bool) {
		return expr.thenBranch.Interpret(environment)
	} else if expr.elseBranch != nil {
		return expr.elseBranch.Interpret(environment)
	}

	return nil, nil
}

func (expr WhileExpr) Interpret(environment Environment) (any, error) {
  condVal, err := expr.condition.Interpret(environment)
  if err != nil {
    return nil, err
  }
  if _, ok := condVal.(bool); !ok {
    return nil, fmt.Errorf("Expected Type bool, got %T", condVal)
  }

	for condVal.(bool) {
		_, err := expr.loopBranch.Interpret(environment)
    if err != nil {
      return nil, err
    }

    condVal, err = expr.condition.Interpret(environment)
    if err != nil {
      return nil, err
    }
    if _, ok := condVal.(bool); !ok {
      return nil, fmt.Errorf("Expected Type bool, got %T", condVal)
    }
	}

	return nil, nil
}

func (expr CallExpr) Interpret(environment Environment) (any, error) {
	f, err := environment.findVar(expr.f.String())
	if err != nil {
		return nil, err
	}

	function, ok := f.(Function)
	if !ok {
		return nil, fmt.Errorf("Invalid Function %v", function)
	}
	args := []any{}

	for _, arg := range expr.args {
    a, err := arg.Interpret(environment)
    if err != nil {
      return nil, err
    }
		args = append(args, a)
	}

	value, err := function.Call(environment, args)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (expr FnDeclExpr) Interpret(environment Environment) (any, error) {
  f := Function {
    Arity: len(expr.args),
    Call: func(env Environment, args []any) (any, error) {
      env.push(make(map[string]any))

      if len(expr.args) != len(args) {
        return nil, fmt.Errorf("[ERROR] Arity does not match at Function %v at Line %v", expr.name.Value, expr.name.Line)
      }

      for i := range args {
        env.declareVar(expr.args[i].Value.(string), args[i])
      }

      result, err := expr.program.Interpret(env)
      env.pop()
      if err != nil {
        return nil, err
      }

      return result, nil
    },
  }

  environment.declareVar(expr.name.Value.(string), f)

  return f, nil
}
