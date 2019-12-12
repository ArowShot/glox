package main

import (
	"bytes"
	"fmt"

	"github.com/arowshot/glox"
)

type AstPrinter struct {
}

func (ast AstPrinter) VisitBinaryExpr(expr glox.Binary) interface{} {
	return ast.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}
func (ast AstPrinter) VisitGroupingExpr(expr glox.Grouping) interface{} {
	return ast.parenthesize("group", expr.Expression)
}
func (ast AstPrinter) VisitLiteralExpr(expr glox.Literal) interface{} {
	return fmt.Sprintf("%v", expr.Value)
}
func (ast AstPrinter) VisitUnaryExpr(expr glox.Unary) interface{} {
	return ast.parenthesize(expr.Operator.Lexeme, expr.Right)
}
func (ast AstPrinter) VisitVariableExpr(expr glox.Variable) interface{} {
	return fmt.Sprintf("(var %v)", expr.Name.Lexeme)
}
func (ast AstPrinter) VisitAssignExpr(expr glox.Assign) interface{} {
	return ast.parenthesize(fmt.Sprintf("assign %v", expr.Name.Lexeme), expr)
}
func (ast *AstPrinter) visitor() *glox.ExprVisitor {
	var v glox.ExprVisitor = *ast
	return &v
}
func (ast *AstPrinter) print(expr glox.Expr) string {
	return fmt.Sprintf("%v", expr.Accept(ast.visitor()))
}

func (ast AstPrinter) parenthesize(name string, exprs ...glox.Expr) string {
	writer := bytes.NewBufferString("")

	fmt.Fprintf(writer, "(%s", name)
	for _, expr := range exprs {
		fmt.Fprintf(writer, " %s", expr.Accept(ast.visitor()))
	}
	fmt.Fprint(writer, ")")

	return writer.String()
}
