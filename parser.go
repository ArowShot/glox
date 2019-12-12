package glox

import "fmt"

type Parser struct {
	Tokens  []Token
	current int
}

func (parser *Parser) ParseExpr() Expr {
	return parser.readExpression()
}

func (parser *Parser) Parse() []Stmt {
	var stmts []Stmt

	for !parser.atEnd() {
		stmts = append(stmts, parser.readDeclaration())
	}

	return stmts
}

func (parser *Parser) readDeclaration() Stmt {
	if parser.match(VAR) {
		return parser.readVarDeclaration()
	}

	return parser.readStatement()
}

func (parser *Parser) readVarDeclaration() Stmt {
	name, _ := parser.consume(IDENTIFIER, "Expected variable name.")

	var expr Expr
	if parser.match(EQUAL) {
		expr = parser.readExpression()
	}

	parser.consume(SEMICOLON, "Expected ';' after variable declaration")

	return Var{name, expr}
}

func (parser *Parser) readStatement() Stmt {
	if parser.match(PRINT) {
		return parser.readPrintStatement()
	}

	if parser.match(WHILE) {
		return parser.readWhileStatement()
	}

	if parser.match(LEFT_BRACE) {
		return Block{parser.readBlock()}
	}

	return parser.readExpressionStatement()
}

func (parser *Parser) readWhileStatement() Stmt {
	parser.consume(LEFT_PAREN, "expected '(' after while")
	condition := parser.readExpression()
	parser.consume(RIGHT_PAREN, "expected ')' after condition")
	body := parser.readStatement()

	return While{condition, body}
}

func (parser *Parser) readBlock() []Stmt {
	var stmts []Stmt

	for !parser.check(RIGHT_BRACE) && !parser.atEnd() {
		stmts = append(stmts, parser.readDeclaration())
	}

	parser.consume(RIGHT_BRACE, "expected '}' after block")

	return stmts
}


func (parser *Parser) readPrintStatement() Stmt {
	value := parser.readExpression()
	parser.consume(SEMICOLON, "Expected ';' after value.")
	return Print{
		Expression: value,
	}
}

func (parser *Parser) readExpressionStatement() Stmt {
	value := parser.readExpression()
	parser.consume(SEMICOLON, "Expected ';' after expression.")
	return Expression{
		Expression: value,
	}
}

func (parser *Parser) readExpression() Expr {
	return parser.readAssignment()
}

func (parser *Parser) readAssignment() Expr {
	expr := parser.readEquality()

	if parser.match(EQUAL) {
		//equals := parser.previous()
		value := parser.readAssignment()

		if v, isVar := expr.(Variable); isVar {
			name := v.Name
			return Assign{name, value}
		}
	}

	return expr
}

func (parser *Parser) readEquality() Expr {
	expr := parser.readComparison()

	for parser.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := parser.previous()
		right := parser.readComparison()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (parser *Parser) readComparison() Expr {
	expr := parser.readAddition()

	for parser.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.previous()
		right := parser.readAddition()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (parser *Parser) readAddition() Expr {
	expr := parser.readMultiplication()

	for parser.match(PLUS, MINUS) {
		operator := parser.previous()
		right := parser.readMultiplication()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (parser *Parser) readMultiplication() Expr {
	expr := parser.readExponent()

	for parser.match(SLASH, STAR) {
		operator := parser.previous()
		right := parser.readExponent()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (parser *Parser) readExponent() Expr {
	expr := parser.readUnary()

	for parser.match(STARSTAR) {
		operator := parser.previous()
		right := parser.readUnary()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (parser *Parser) readUnary() Expr {
	if parser.match(BANG, MINUS) {
		operator := parser.previous()
		right := parser.readUnary()
		return Unary{operator, right}
	}

	return parser.readPrimary()
}

func (parser *Parser) readPrimary() Expr {
	if parser.match(FALSE) {
		return Literal{false}
	}
	if parser.match(TRUE) {
		return Literal{true}
	}
	if parser.match(NIL) {
		return Literal{nil}
	}

	if parser.match(NUMBER, STRING) {
		return Literal{parser.previous().Literal}
	}

	if parser.match(IDENTIFIER) {
		return Variable{parser.previous()}
	}

	if parser.match(LEFT_PAREN) {
		expr := parser.readExpression()
		parser.consume(RIGHT_PAREN, "Unclosed '('.")
		return Grouping{expr}
	}

	return nil
}

func (parser *Parser) consume(tokenType int, message string) (Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	}

	err := fmt.Errorf("error on line %v at\"%s\": %s", parser.peek().Line, parser.peek().Lexeme, message)
	fmt.Println(err)

	return Token{}, err
}

func (parser *Parser) match(tokenTypes ...int) bool {
	for _, tokenType := range tokenTypes {
		if parser.check(tokenType) {
			parser.advance()
			return true
		}
	}

	return false
}

func (parser *Parser) check(tokenType int) bool {
	if parser.atEnd() {
		return false
	}
	return parser.peek().TokenType == tokenType
}

func (parser *Parser) advance() Token {
	if !parser.atEnd() {
		parser.current++
	}
	return parser.previous()
}

func (parser *Parser) previous() Token {
	return parser.Tokens[parser.current-1]
}

func (parser *Parser) peek() Token {
	return parser.Tokens[parser.current]
}

func (parser *Parser) atEnd() bool {
	return parser.peek().TokenType == EOF
}
