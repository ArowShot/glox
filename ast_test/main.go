package main

import (
	"bytes"
	"fmt"

	"github.com/arowshot/glox"
)

type AstPrinter struct {
}

func (ast AstPrinter) VisitBinaryExpr(expr glox.Binary) interface{} {
	return ast.parenthesize(expr.Operator.Lexeme, []glox.Expr{expr.Left, expr.Right})
}
func (ast AstPrinter) VisitGroupingExpr(expr glox.Grouping) interface{} {
	return ast.parenthesize("group", []glox.Expr{expr.Expression})
}
func (ast AstPrinter) VisitLiteralExpr(expr glox.Literal) interface{} {
	return fmt.Sprintf("%v", expr.Value)
}
func (ast AstPrinter) VisitUnaryExpr(expr glox.Unary) interface{} {
	return ast.parenthesize(expr.Operator.Lexeme, []glox.Expr{expr.Right})
}
func (ast *AstPrinter) visitor() *glox.ExprVisitor {
	var v glox.ExprVisitor = *ast
	return &v
}
func (ast *AstPrinter) print(expr glox.Expr) string {
	return fmt.Sprintf("%v", expr.Accept(ast.visitor()))
}

func (ast AstPrinter) parenthesize(name string, exprs []glox.Expr) string {
	writer := bytes.NewBufferString("")

	fmt.Fprintf(writer, "(%s", name)
	for _, expr := range exprs {
		fmt.Fprintf(writer, " %s", expr.Accept(ast.visitor()))
	}
	fmt.Fprint(writer, ")")

	return writer.String()
}

func main() {
	expr := glox.Binary{
		Left: glox.Unary{
			Operator: glox.Token{
				Lexeme:    "-",
				TokenType: glox.MINUS,
			},
			Right: glox.Literal{
				Value: 3.0,
			},
		},
		Operator: glox.Token{
			Lexeme:    "*",
			TokenType: glox.SLASH,
		},
		Right: glox.Grouping{
			Expression: glox.Binary{
				Left: glox.Unary{
					Operator: glox.Token{
						Lexeme:    "-",
						TokenType: glox.MINUS,
					},
					Right: glox.Literal{
						Value: 3.0,
					},
				},
				Operator: glox.Token{
					Lexeme:    "*",
					TokenType: glox.SLASH,
				},
				Right: glox.Grouping{
					Expression: glox.Literal{
						Value: 47.56,
					},
				},
			},
		},
	}

	astPrinter := AstPrinter{}

	fmt.Println(astPrinter.print(expr))

	expr2 := glox.Binary{
		Left: glox.Literal{
			Value: 3.0,
		},
		Operator: glox.Token{
			Lexeme:    "**",
			TokenType: glox.STARSTAR,
		},
		Right: glox.Literal{
			Value: 3.0,
		},
	}

	interpreter := glox.Interpreter{}

	interpreter.Interpret(expr2)
}
