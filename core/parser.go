package core

import (
	"fmt"
	"slices"
)

type Parser struct {
	tokens  []Token
	current int
}

func CreateParser(tokens []Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (parser *Parser) Parse() ([]Expr, error) {
	var program []Expr = []Expr{}

	for !parser.isAtEnd() {
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}

		program = append(program, expr)
	}

	return program, nil
}

func (parser *Parser) expression() (Expr, error) {
	return parser.controlFlow()
}

func (parser *Parser) controlFlow() (Expr, error) {
	if parser.match(IF) {
		return parser.ifStmt()
	}
	if parser.match(WHILE) {
		return parser.whiteStmt()
	}
	if parser.match(FOR) {
		return parser.forStmt()
	}
  if parser.match(FUNCTION) {
    return parser.fnDeclStmt()
  }

	return parser.block()
}

func (parser *Parser) ifStmt() (Expr, error) {
	_, err := parser.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := parser.expression()
	if err != nil {
		return nil, err
	}

	_, err = parser.consume(RIGHT_PAREN, "Expected ) after condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := parser.expression()
	if err != nil {
		return nil, err
	}

	var elseBranch Expr
	if parser.match(ELSE) {
		elseBranch, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	return IfExpr{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}, nil
}

func (parser *Parser) whiteStmt() (Expr, error) {
	_, err := parser.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := parser.expression()
	if err != nil {
		return nil, err
	}

	_, err = parser.consume(RIGHT_PAREN, "Expected ) after condition")
	if err != nil {
		return nil, err
	}

	loopBranch, err := parser.expression()
	if err != nil {
		return nil, err
	}

	return WhileExpr{
		condition:  condition,
		loopBranch: BlockExpr{ 
      program: []Expr{loopBranch},
    },
	}, nil
}

func (parser *Parser) forStmt() (Expr, error) {
	_, err := parser.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	initializer, err := parser.expression()
	if err != nil {
		return nil, err
	}
	_, err = parser.consume(SEMI, "Expected ; after condition")
	if err != nil {
		return nil, err
	}

	condition, err := parser.expression()
	if err != nil {
		return nil, err
	}
	_, err = parser.consume(SEMI, "Expected ; after condition")
	if err != nil {
		return nil, err
	}

	increment, err := parser.expression()
	if err != nil {
		return nil, err
	}

	_, err = parser.consume(RIGHT_PAREN, "Expected ) after condition")
	if err != nil {
		return nil, err
	}

	loopBranchContent, err := parser.expression()
	if err != nil {
		return nil, err
	}

	loopBranch := BlockExpr{
		program: []Expr{loopBranchContent, increment},
	}

	return BlockExpr{
		program: append([]Expr{initializer}, WhileExpr{
			condition:  condition,
			loopBranch: loopBranch,
		}),
	}, nil
}

func (parser *Parser) fnDeclStmt() (Expr, error) {
  identifier, err := parser.consume(IDENTIFIER, "Expected Identifier after fn")
  if err != nil {
    return nil, err
  }

  _, err = parser.consume(LEFT_PAREN, "Expected ( after fnDecl")
  if err != nil {
    return nil, err
  }

  args, err := parser.finishArgs()
  if err != nil {
    return nil, err
  }

  program, err := parser.expression()
  if err != nil {
    return nil, err
  }

  return FnDeclExpr {
    name: identifier,
    args: args,
    program: program,
  }, nil
}

func (parser *Parser) finishArgs() ([]Token, error) {
	var args []Token

	if !parser.check(RIGHT_PAREN) {
		if parser.isAtEnd() {
			return nil, fmt.Errorf("Expected ')' after args")
		}

		arg, err := parser.consume(IDENTIFIER, "Expected Identifier after (")
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		for parser.match(COMMA) {
			arg, err := parser.consume(IDENTIFIER, "Expcted Identifier after ,")
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
	}

	_, err := parser.consume(RIGHT_PAREN, "Expected ')' after args")
  if err != nil {
    return nil, err
  }

	return args, nil
}

func (parser *Parser) block() (Expr, error) {
	if parser.match(LEFT_BRACE) {
		program := []Expr{}
		for !parser.match(RIGHT_BRACE) {
			if parser.isAtEnd() {
				return nil, fmt.Errorf("Expected } after block")
			}

			expr, err := parser.expression()
			if err != nil {
				return nil, err
			}

			program = append(program, expr)
		}

		return BlockExpr{
			program: program,
		}, nil
	}

	return parser.assignment()
}

func (parser *Parser) assignment() (Expr, error) {
	if parser.match(IDENTIFIER) {
		name := parser.tokens[parser.current-1]
		for parser.match(EQUAL, COLON_EQ, PLUS_EQ, MINUS_EQ, STAR_EQ, SLASH_EQ, PERCENT_EQ) {
			operator := parser.tokens[parser.current-1]
			expr, err := parser.expression()
			if err != nil {
				return nil, err
			}

			switch operator.Type {
			case PLUS_EQ:
				return parser.assignDesugared(name, PLUS, expr, operator.Line), nil
			case MINUS_EQ:
				return parser.assignDesugared(name, MINUS, expr, operator.Line), nil
			case STAR_EQ:
				return parser.assignDesugared(name, STAR, expr, operator.Line), nil
			case SLASH_EQ:
				return parser.assignDesugared(name, SLASH, expr, operator.Line), nil
			case PERCENT_EQ:
				return parser.assignDesugared(name, PERCENT, expr, operator.Line), nil
			}

			return AssignExpr{
				name:     name,
				expr:     expr,
				operator: operator,
			}, nil
		}
		parser.current--
	}
	return parser.logicalOr()
}

func (parser *Parser) assignDesugared(name Token, operator TokenType, expr Expr, line int) AssignExpr {
	return AssignExpr{
		name: name,
		expr: BinaryExpr{
			operator: Token{
				Type:  operator,
				Value: "",
				Line:  line,
			},
			leftExpr: LiteralExpr{
				value: Token{
					Type:  IDENTIFIER,
					Value: name.Value,
					Line:  line,
				},
			},
			rightExpr: expr,
		},
		operator: Token{
			Type: EQUAL, Value: "", Line: line,
		},
	}
}

func (parser *Parser) logicalOr() (Expr, error) {
	expr, err := parser.logicalAnd()
	if err != nil {
		return nil, err
	}

	for parser.match(BAR_BAR) {
		operator := parser.tokens[parser.current-1]
		rightExpr, err := parser.logicalAnd()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr, nil
}

func (parser *Parser) logicalAnd() (Expr, error) {
	expr, err := parser.equality()
	if err != nil {
		return nil, err
	}

	for parser.match(AMP_AMP) {
		operator := parser.tokens[parser.current-1]
		rightExpr, err := parser.equality()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr, nil
}

func (parser *Parser) equality() (Expr, error) {
	expr, err := parser.comparison()
	if err != nil {
		return nil, err
	}

	for parser.match(EQUAL_EQ, BANG_EQ) {
		operator := parser.tokens[parser.current-1]
		rightExpr, err := parser.comparison()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr, nil
}

func (parser *Parser) comparison() (Expr, error) {
	expr, err := parser.term()
	if err != nil {
		return nil, err
	}

	for parser.match(GREATER, GREATER_EQ, LESS, LESS_EQ) {
		operator := parser.tokens[parser.current-1]
		rightExpr, err := parser.comparison()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr, nil
}

func (parser *Parser) term() (Expr, error) {
	expr, err := parser.factor()
	if err != nil {
		return nil, err
	}

	for parser.match(PLUS, MINUS) {
		operator := parser.tokens[parser.current-1]
		rightExpr, err := parser.factor()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr, nil
}

func (parser *Parser) factor() (Expr, error) {
	expr, err := parser.unary()
	if err != nil {
		return nil, err
	}

	for parser.match(STAR, SLASH, PERCENT) {
		operator := parser.tokens[parser.current-1]
		rightExpr, err := parser.unary()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			operator:  operator,
			rightExpr: rightExpr,
			leftExpr:  expr,
		}
	}

	return expr, nil
}

func (parser *Parser) unary() (Expr, error) {
	if parser.match(BANG, MINUS, HASHTAG) {
    operator := parser.tokens[parser.current-1]
		expr, err := parser.index()
		if err != nil {
			return nil, err
		}
		return UnaryExpr{
			operator: operator,
			expr:     expr,
		}, nil
	}

	expr, err := parser.index()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (parser *Parser) index() (Expr, error) {
  expr, err := parser.call()
  if err != nil {
    return nil, err
  }

  if parser.match(LEFT_BRACKET) {
    arg, err := parser.expression()
    if err != nil {
      return nil, err
    }
    _, err = parser.consume(RIGHT_BRACKET, "Expected ']' after index notation")
    if err != nil {
      return nil, err
    }

    return IndexExpr{
      value: expr, 
      index: arg,
    }, nil
  }

  return expr, nil
}

func (parser *Parser) call() (Expr, error) {
	expr, err := parser.primary()
	if err != nil {
		return nil, err
	}

	if parser.match(LEFT_PAREN) {
		args, err := parser.finishCall()
		if err != nil {
			return nil, err
		}

		return CallExpr{
			f:    expr,
			args: args,
		}, nil
	}

	return expr, nil
}

func (parser *Parser) finishCall() ([]Expr, error) {
	var args []Expr

	if !parser.check(RIGHT_PAREN) {
		if parser.isAtEnd() {
			return nil, fmt.Errorf("Expected ')' after args")
		}

		arg, err := parser.expression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		for parser.match(COMMA) {
			arg, err := parser.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
	}

	_, err := parser.consume(RIGHT_PAREN, "Expected ')' after args")
  if err != nil {
    return nil, err
  }

	return args, nil
}

func (parser *Parser) primary() (Expr, error) {
	if parser.match(IDENTIFIER, NUMBER, STRING, TRUE, FALSE) {
		return LiteralExpr{
			value: parser.tokens[parser.current-1],
		}, nil
	}

  if parser.match(LEFT_BRACKET) {
    var values []Expr
    if !parser.check(RIGHT_BRACKET) {
      if parser.isAtEnd() {
        return nil, fmt.Errorf("Expected ']' after array initializer")
      }

      val, err := parser.expression()
      if err != nil {
        return nil, err
      }
      values = append(values, val)

      for parser.match(COMMA) {
        val, err := parser.expression()
        if err != nil {
          return nil, err
        }
        values = append(values, val)
      }
    }
    
    _, err := parser.consume(RIGHT_BRACKET, "Expected ']' after array initializer")
    if err != nil {
      return nil, err
    }

    return ArrayInitExpr {
      values: values,
    }, nil
  }

	if parser.match(LEFT_PAREN) {
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}
		_, err = parser.consume(RIGHT_PAREN, "Expected ')' after expression")
		if err != nil {
			return nil, err
		}

		return GroupingExpr{
			expr: expr,
		}, nil
	}

	return nil, fmt.Errorf("[ERROR] Syntax Error at Line %d\n", parser.tokens[parser.current].Line)
}

func (parser *Parser) match(tokenTypes ...TokenType) bool {
	if slices.Contains(tokenTypes, parser.tokens[parser.current].Type) {
		parser.advance()
		return true
	}

	return false
}

func (parser *Parser) advance() Token {
	result := parser.tokens[parser.current]

	if !parser.isAtEnd() {
		parser.current++
	}

	return result
}

func (parser *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	}

	return parser.advance(), fmt.Errorf("[ERROR] %s at Line %d", message, parser.tokens[parser.current].Line)
}

func (parser Parser) check(tokenType TokenType) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.tokens[parser.current].Type == tokenType
}

func (parser Parser) isAtEnd() bool {
	return parser.current == len(parser.tokens)-1
}
