package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	exprAst := defineAst("Expr", []string{
		"Assign : Name Token, Value Expr",
		"Binary : Left Expr, Operator Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal : Value interface{}",
		"Unary : Operator Token, Right Expr",
		"Variable : Name Token",
	})

	file, err := os.Create("expressions.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.WriteString(file, exprAst)
	if err != nil {
		panic(err)
	}

	stmtAst := defineAst("Stmt", []string{
		"Block : Statements []Stmt",
		"Expression : Expression Expr",
		"Print : Expression Expr",
		"Var : Name Token, Initializer Expr",
		"While : Condition Expr, Body Stmt",
	})

	file, err = os.Create("statements.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.WriteString(file, stmtAst)
	if err != nil {
		panic(err)
	}
}

func defineAst(baseType string, types []string) string {
	writer := bytes.NewBufferString("")

	fmt.Fprintln(writer, "package glox")
	fmt.Fprintln(writer, "")
	fmt.Fprintf(writer, "type %s interface {\n", baseType)
	fmt.Fprintf(writer, "    Accept(*%sVisitor) interface{}\n", baseType)
	fmt.Fprintln(writer, "}")
	fmt.Fprintln(writer, "")

	defineVisitor(writer, baseType, types)

	for _, t := range types {
		exprName := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(writer, baseType, exprName, fields)
	}

	return writer.String()
}

func defineType(writer *bytes.Buffer, baseName string, exprName string, fields string) {
	fmt.Fprintf(writer, "type %s struct {\n", exprName)
	fieldList := strings.Split(fields, ", ")
	for _, field := range fieldList {
		fmt.Fprintf(writer, "    %s\n", field)
	}
	fmt.Fprintln(writer, "}")

	fmt.Fprintf(writer, "func (me %s) Accept(visitor *%sVisitor) interface{} {\n", exprName, baseName)
	fmt.Fprintf(writer, "    v := *visitor\n")
	fmt.Fprintf(writer, "    return v.Visit%s%s(me)\n", exprName, baseName)
	fmt.Fprintln(writer, "}")

	fmt.Fprintln(writer, "")
}

func defineVisitor(writer *bytes.Buffer, baseType string, types []string) {
	fmt.Fprintf(writer, "type %sVisitor interface {\n", baseType)
	for _, t := range types {
		typeName := strings.TrimSpace(strings.Split(t, ":")[0])
		fmt.Fprintf(writer, "    Visit%s%s(%s) interface{}\n", typeName, baseType, typeName)
	}
	fmt.Fprintln(writer, "}")
	fmt.Fprintln(writer, "")
}
