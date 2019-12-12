package glox

import (
	"fmt"
	"strconv"
)

// Keywords are the keywords defined in the language
var Keywords = map[string]int{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

// Scanner is used to generate tokens from a source text
type Scanner struct {
	Source         string
	Tokens         []Token
	start, current int
	line           int
}

func (sc *Scanner) atEnd() bool {
	return sc.current >= len([]rune(sc.Source))
}

func (sc *Scanner) advance() rune {
	sc.current++
	return []rune(sc.Source)[sc.current-1]
}

func (sc *Scanner) peek() rune {
	if sc.atEnd() {
		return '\x00'
	}
	return []rune(sc.Source)[sc.current]
}

func (sc *Scanner) peekNext() rune {
	if sc.current+1 >= len([]rune(sc.Source)) {
		return '\x00'
	}
	return []rune(sc.Source)[sc.current+1]
}

func (sc *Scanner) match(c rune) bool {
	if sc.atEnd() {
		return false
	}

	if []rune(sc.Source)[sc.current] != c {
		return false
	}

	sc.current++
	return true
}

func (sc *Scanner) addToken(token int) {
	sc.addLiteralToken(token, nil)
}

func (sc *Scanner) addLiteralToken(token int, literal interface{}) {
	text := string([]rune(sc.Source)[sc.start:sc.current])
	sc.Tokens = append(sc.Tokens, Token{
		TokenType: token,
		Lexeme:    text,
		Literal:   literal,
		Line:      sc.line,
	})
}

func (sc *Scanner) scanStr() {
	for sc.peek() != '"' && !sc.atEnd() { // Keep reading characters until a " is found or the end is reached
		if sc.peek() == '\n' { // If we encounter a newline
			sc.line++ // increase the line count
		}
		sc.advance() // Advance the scanner to the next character
	}

	if sc.atEnd() { // If we're at the end of the source then there was never a closing "
		fmt.Printf("Unterminated string")
		return
	}

	sc.advance() // Advance last the final "

	value := []rune(sc.Source)[sc.start+1 : sc.current-1] // Get the value of the string by reading from the source
	sc.addLiteralToken(STRING, string(value)) // Add the string token to the token list
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '_')
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (sc *Scanner) scanNum() {
	for isDigit(sc.peek()) {
		sc.advance()
	}

	if sc.peek() == '.' && isDigit(sc.peekNext()) {
		sc.advance()

		for isDigit(sc.peek()) {
			sc.advance()
		}
	}

	num, _ := strconv.ParseFloat(string([]rune(sc.Source)[sc.start:sc.current]), 64)

	sc.addLiteralToken(NUMBER, num)
}

func (sc *Scanner) scanIdentifier() {
	for isAlphaNumeric(sc.peek()) {
		sc.advance()
	}

	text := string([]rune(sc.Source[sc.start:sc.current]))

	tokenType := Keywords[text]
	if tokenType == 0 {
		tokenType = IDENTIFIER
	}

	sc.addToken(tokenType)
}

func (sc *Scanner) scanToken() {
	c := sc.advance()
	switch c {
	case '(':
		sc.addToken(LEFT_PAREN)
	case ')':
		sc.addToken(RIGHT_PAREN)
	case '{':
		sc.addToken(LEFT_BRACE)
	case '}':
		sc.addToken(RIGHT_BRACE)
	case ',':
		sc.addToken(COMMA)
	case '.':
		sc.addToken(DOT)
	case '-':
		sc.addToken(MINUS)
	case '+':
		sc.addToken(PLUS)
	case ';':
		sc.addToken(SEMICOLON)
	case '*':
		if sc.match('*') {
			sc.addToken(STARSTAR)
		} else {
			sc.addToken(STAR)
		}
	case '!':
		if sc.match('=') {
			sc.addToken(BANG_EQUAL)
		} else {
			sc.addToken(BANG)
		}
	case '=':
		if sc.match('=') {
			sc.addToken(EQUAL_EQUAL)
		} else {
			sc.addToken(EQUAL)
		}
	case '<':
		if sc.match('=') {
			sc.addToken(LESS_EQUAL)
		} else {
			sc.addToken(LESS)
		}
	case '>':
		if sc.match('=') {
			sc.addToken(GREATER_EQUAL)
		} else {
			sc.addToken(GREATER)
		}
	case '/':
		if sc.match('/') {
			for sc.peek() != '\n' && !sc.atEnd() {
				sc.advance()
			}
		} else {
			sc.addToken(SLASH)
		}
	case ' ':
		fallthrough
	case '\r':
		fallthrough
	case '\t':
	case '\n':
		sc.line++
	case '"':
		sc.scanStr()
	default:
		if isDigit(c) {
			sc.scanNum()
		} else if isAlphaNumeric(c) {
			sc.scanIdentifier()
		} else {
			fmt.Printf("Unexpected character on line %v\n", sc.line)
		}
	}
}

// ScanTokens will scan the soruce code for tokens
func (sc *Scanner) ScanTokens() []Token {
	for !sc.atEnd() { // Keep reading every character until at the end of the file
		sc.start = sc.current // The token starts at the current position
		sc.scanToken() // Scan a token from the source
	}
	
	sc.addToken(EOF) // Add en end of file token to the end

	return sc.Tokens // Return the list of tokens
}

