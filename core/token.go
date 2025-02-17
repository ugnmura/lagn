package core

import (
	"fmt"
)

type TokenType int
type TokenValue interface{}

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE

	COMMA
	DOT

	MINUS
	MINUS_MINUS
	MINUS_EQ
	PLUS
	PLUS_PLUS
	PLUS_EQ
	SEMI
	SLASH
	SLASH_EQ
	STAR
	STAR_EQ

	BIN_AND
	BIN_AND_EQ
	BIN_OR
	BIN_OR_EQ
	BIN_XOR
	BIN_XOR_EQ

	BANG
	BANG_EQ
	EQUAL
	EQUAL_EQ
	GREATER
	GREATER_EQ
	LESS
	LESS_EQ

	IDENTIFIER
	STRING
	NUMBER

	AND
	OR
	FOR
	WHILE
	IF
	ELSE
	RETURN

	TRUE
	FALSE

	UNKOWN

	EOF
)

var tokenTypeNames = [...]string{
	"LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE",
	"COMMA", "DOT",
	"MINUS", "MINUS_MINUS", "MINUS_EQ", "PLUS", "PLUS_PLUS", "PLUS_EQ",
	"SEMI", "SLASH", "SLASH_EQ", "STAR", "STAR_EQ",
	"BIN_AND", "BIN_AND_EQ", "BIN_OR", "BIN_OR_EQ", "BIN_XOR", "BIN_XOR_EQ",
	"BANG", "BANG_EQ", "EQUAL", "EQUAL_EQ", "GREATER", "GREATER_EQ", "LESS", "LESS_EQ",
	"IDENTIFIER", "STRING", "NUMBER",
	"AND", "OR", "FOR", "WHILE", "IF", "ELSE", "RETURN",
	"TRUE", "FALSE",
	"UNKOWN",
	"EOF",
}

func (tt TokenType) String() string {
	return tokenTypeNames[tt]
}

func init() {
	if len(tokenTypeNames)-1 != int(EOF) {
		fmt.Println("[ERROR] TokenTypes are not updated correctly")
	}
}

type Token struct {
	Type  TokenType
	Value TokenValue
	Line  int
}

func (token Token) String() string {
	if token.Value == nil {
		return token.Type.String()
	}
	if token.Type == STRING {
		return fmt.Sprintf("%q", token.Value.(string))
	}
	if token.Type == NUMBER {
		return fmt.Sprintf("%v", token.Value.(float64))
	}
	return token.Value.(string)
}
