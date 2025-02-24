package core

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

var KEYWORDS map[string]TokenType

func init() {
	KEYWORDS = make(map[string]TokenType)
	KEYWORDS["for"] = FOR
	KEYWORDS["while"] = WHILE
	KEYWORDS["else"] = ELSE
	KEYWORDS["if"] = IF
	KEYWORDS["return"] = RETURN
	KEYWORDS["true"] = TRUE
	KEYWORDS["false"] = FALSE
}

type Scanner struct {
	Source  []rune
	Tokens  []Token
	Start   int
	Current int
	Line    int
	Column  int
}

func CreateScanner(source string) Scanner {
	return Scanner{
		Source:  []rune(source),
		Start:   0,
		Current: 0,
		Line:    1,
		Column:  1,
	}
}

func (scanner *Scanner) AddToken(token TokenType) {
	scanner.Tokens = append(scanner.Tokens, Token{
		Type: token,
		Line: scanner.Line,
	})
}

func (scanner *Scanner) AddTokenWithValue(token TokenType, value TokenValue) {
	scanner.Tokens = append(scanner.Tokens, Token{
		Type:  token,
		Line:  scanner.Line,
		Value: value,
	})
}

func (scanner *Scanner) Advance() rune {
	result := scanner.Source[scanner.Current]
	scanner.Current++
	scanner.Column++
	return result
}

func (scanner *Scanner) ScanToken() {
	c := scanner.Advance()
	switch c {
	case rune('('):
		scanner.AddToken(LEFT_PAREN)
	case rune(')'):
		scanner.AddToken(RIGHT_PAREN)
	case rune('{'):
		scanner.AddToken(LEFT_BRACE)
	case rune('}'):
		scanner.AddToken(RIGHT_BRACE)
	case rune('+'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(PLUS_EQ)
		} else if scanner.PeekCurrent() == rune('+') {
			scanner.Advance()
			scanner.AddToken(PLUS_PLUS)
		} else {
			scanner.AddToken(PLUS)
		}
	case rune('-'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(MINUS_EQ)
		} else if scanner.PeekCurrent() == rune('-') {
			scanner.Advance()
			scanner.AddToken(MINUS_MINUS)
		} else {
			scanner.AddToken(MINUS)
		}
	case rune('*'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(STAR_EQ)
		} else {
			scanner.AddToken(STAR)
		}
	case rune('/'):
		if scanner.PeekCurrent() == rune('/') {
			scanner.Advance()
			for !scanner.CurrentAtEnd() && scanner.PeekCurrent() != '\n' {
				scanner.Advance()
			}
		} else if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(SLASH_EQ)
		} else {
			scanner.AddToken(SLASH)
		}
	case rune('='):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(EQUAL_EQ)
		} else {
			scanner.AddToken(EQUAL)
		}
	case rune('>'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(GREATER_EQ)
		} else {
			scanner.AddToken(GREATER)
		}
	case rune('<'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(LESS_EQ)
		} else {
			scanner.AddToken(LESS)
		}
	case rune('!'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(BANG_EQ)
		} else {
			scanner.AddToken(BANG)
		}
	case rune(':'):
		if scanner.PeekCurrent() == rune('=') {
			scanner.Advance()
			scanner.AddToken(COLON_EQ)
		} else {
			scanner.AddToken(COLON)
		}
	case rune('.'):
		scanner.AddToken(DOT)
	case rune(','):
		scanner.AddToken(COMMA)
	case rune(';'):
		scanner.AddToken(SEMI)
	case rune('"'):
		scanner.ScanString()
	case rune(' '):
	case rune('\t'):
	case rune('\r'):
	case rune('\n'):
		scanner.Line++
		scanner.Column = 1
	default:
		if unicode.IsDigit(c) {
			scanner.ScanNumber()
		} else if unicode.IsLetter(c) {
			scanner.ScanIdentifier()
		} else {
			fmt.Printf("[ERROR] Unexpected Character at Line %d, Column %d\n", scanner.Line, scanner.Column)
			os.Exit(1)
		}
	}
}

func (scanner *Scanner) ScanNumber() {
	for unicode.IsDigit(scanner.PeekCurrent()) {
		scanner.Advance()
	}
	if scanner.PeekCurrent() == rune('.') && unicode.IsDigit(scanner.Peek(scanner.Current+1)) {
		scanner.Advance()
		for unicode.IsDigit(scanner.PeekCurrent()) {
			scanner.Advance()
		}
	}

	value, _ := strconv.ParseFloat(string(scanner.Source[scanner.Start:scanner.Current]), 64)
	scanner.AddTokenWithValue(NUMBER, value)
}

func (scanner *Scanner) ScanString() {
	for !scanner.CurrentAtEnd() && scanner.PeekCurrent() != rune('"') {
		scanner.Advance()
	}

	if scanner.CurrentAtEnd() {
		fmt.Printf("[ERROR] Unterminated String at Line %d, Column %d\n", scanner.Line, scanner.Column)
		return
	}

	scanner.Advance()
	scanner.AddTokenWithValue(STRING, string(scanner.Source[scanner.Start+1:scanner.Current-1]))
}

func (scanner *Scanner) ScanIdentifier() {
	for !scanner.CurrentAtEnd() && (unicode.IsLetter(scanner.PeekCurrent()) || unicode.IsDigit(scanner.PeekCurrent())) {
		scanner.Advance()
	}

	text := scanner.Source[scanner.Start:scanner.Current]
	tokenType, ok := KEYWORDS[string(text)]
	if !ok {
		tokenType = IDENTIFIER
	}

	scanner.AddTokenWithValue(tokenType, string(text))
}

func (scanner *Scanner) PeekCurrent() rune {
	return scanner.Peek(scanner.Current)
}

func (scanner *Scanner) Peek(toPeek int) rune {
	if scanner.IsAtEnd(toPeek) {
		return rune(0)
	}
	return scanner.Source[toPeek]
}

func (scanner Scanner) CurrentAtEnd() bool {
	return scanner.IsAtEnd(scanner.Current)
}

func (scanner Scanner) IsAtEnd(toCheck int) bool {
	return toCheck >= len(scanner.Source)
}

func (scanner *Scanner) ScanTokens() {
	for !scanner.CurrentAtEnd() {
		scanner.Start = scanner.Current
		scanner.ScanToken()
	}

	scanner.Tokens = append(scanner.Tokens, Token{
		Type: EOF,
		Line: scanner.Line,
	})
}
