package glox

import "fmt"

// Token types
const (
	_ = iota
	// Single-character tokens

	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	STARSTAR

	// One or two character tokens

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals

	IDENTIFIER
	STRING
	NUMBER

	// Keywords

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	// End of file

	EOF
)

// Token represents a single token of the source code
type Token struct {
	TokenType int
	Lexeme    string
	Literal   interface{}
	Line      int
}

func (token Token) String() string {
	return fmt.Sprintf("%v\t%s\t%v", token.TokenType, token.Lexeme, token.Literal)
}
