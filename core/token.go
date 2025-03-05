package core

import (
	"fmt"
)

type TokenType int
type TokenValue any

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_BRACKET
	RIGHT_BRACKET

	COMMA
	DOT
	HASHTAG

	MINUS
	MINUS_MINUS
	MINUS_EQ
	PLUS
	PLUS_PLUS
	PLUS_EQ
	SEMI
	COLON
	COLON_EQ
	SLASH
	SLASH_EQ
	STAR
	STAR_EQ
	PERCENT
	PERCENT_EQ

	AMP
	AMP_AMP
	AMP_EQ
	AMP_AMP_EQ
	BAR
	BAR_BAR
	BAR_EQ
	BAR_BAR_EQ
	CIRCUM
	CIRCUM_EQ
	CIRCUM_CIRCUM
	CIRCUM_CIRCUM_EQ

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

	FOR
	WHILE
	IF
	ELSE
	RETURN
	FUNCTION

	TRUE
	FALSE

	UNKOWN

	EOF
)

var tokenTypeNames = [...]string{
	"LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE", "LEFT_BRACKET", "RIGHT_BRACKET",
	"COMMA", "DOT", "HASHTAG",
	"MINUS", "MINUS_MINUS", "MINUS_EQ", "PLUS", "PLUS_PLUS", "PLUS_EQ",
	"SEMI", "COLON", "COLON_EQ", "SLASH", "SLASH_EQ", "STAR", "STAR_EQ", "PERCENT", "PERCENT_EQ",
	"AMP", "AMP_AMP", "AMP_EQ", "AMP_AMP_EQ",
	"BAR", "BAR_BAR", "BAR_EQ", "BAR_BAR_EQ",
	"CIRCUM", "CIRCUM_EQ", "CIRCUM_CIRCUM", "CIRCUM_CIRCUM_EQ",
	"BANG", "BANG_EQ", "EQUAL", "EQUAL_EQ", "GREATER", "GREATER_EQ", "LESS", "LESS_EQ",
	"IDENTIFIER", "STRING", "NUMBER",
	"FOR", "WHILE", "IF", "ELSE", "RETURN", "FUNCTION",
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
