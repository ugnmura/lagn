package core

import (
	"fmt"
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
		return parser.ifStatement()
	}
	if parser.match(WHILE) {
		return parser.whileStatement()
	}

	return parser.block()
}

func (parser *Parser) ifStatement() (Expr, error) {
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

	thenBranch, err := parser.block()
	if err != nil {
		return nil, err
	}

	var elseBranch Expr
	if parser.match(ELSE) {
		elseBranch, err = parser.block()
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

func (parser *Parser) whileStatement() (Expr, error) {
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

	loopBranch, err := parser.block()
	if err != nil {
		return nil, err
	}

	return WhileExpr{
		condition:  condition,
		loopBranch: loopBranch,
	}, nil
}

func (parser *Parser) block() (Expr, error) {
	if parser.match(LEFT_BRACE) {
		program := []Expr{}
		for !parser.match(RIGHT_BRACE) {
			if parser.isAtEnd() {
				return nil, fmt.Errorf("Expected } after block")
			}

			expr, err := parser.block()
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
		for parser.match(EQUAL) {
			expr, err := parser.expression()
			if err != nil {
				return nil, err
			}

			return AssignExpr{
				name: name,
				expr: expr,
			}, nil
		}
		parser.current--
	}
	return parser.equality()
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

	for parser.match(STAR, SLASH) {
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
	if parser.match(BANG, MINUS) {
		expr, err := parser.unary()
		if err != nil {
			return nil, err
		}
		return UnaryExpr{
			operator: parser.tokens[parser.current-1],
			expr:     expr,
		}, nil
	}

	expr, err := parser.primary()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (parser *Parser) primary() (Expr, error) {
	if parser.match(IDENTIFIER, NUMBER, STRING, TRUE, FALSE) {
		return LiteralExpr{
			value: parser.tokens[parser.current-1],
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
	for _, tokenType := range tokenTypes {
		if parser.tokens[parser.current].Type == tokenType {
			parser.advance()
			return true
		}
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
