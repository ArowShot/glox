package glox

import (
	"fmt"
	"math"
)

type Interpreter struct {
	Env *Environment
}

func (intr Interpreter) VisitBinaryExpr(expr Binary) interface{} {
	left := intr.eval(expr.Left)
	right := intr.eval(expr.Right)

	switch expr.Operator.TokenType {
	case GREATER:
		l := left.(float64)
		r := right.(float64)
		return l > r
	case GREATER_EQUAL:
		l := left.(float64)
		r := right.(float64)
		return l >= r
	case LESS:
		l := left.(float64)
		r := right.(float64)
		return l < r
	case LESS_EQUAL:
		l := left.(float64)
		r := right.(float64)
		return l <= r
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case MINUS:
		l := left.(float64)
		r := right.(float64)
		return l - r
	case STAR:
		l := left.(float64)
		r := right.(float64)
		return l * r
	case STARSTAR:
		l := left.(float64)
		r := right.(float64)
		return math.Pow(l, r)
	case SLASH:
		l := left.(float64)
		r := right.(float64)
		return l / r
	case PLUS:
		lfloat, lisfloat := left.(float64)
		rfloat, risfloat := right.(float64)

		lstring, lisstring := left.(string)
		rstring, risstring := right.(string)

		if lisfloat && risfloat {
			return lfloat + rfloat
		}

		if lisstring && risstring {
			return lstring + rstring
		}
	}

	return nil
}

func (intr Interpreter) VisitGroupingExpr(expr Grouping) interface{} {
	return intr.eval(expr.Expression)
}

func (intr Interpreter) VisitLiteralExpr(expr Literal) interface{} {
	return expr.Value
}

func (intr Interpreter) VisitUnaryExpr(expr Unary) interface{} {
	right := intr.eval(expr.Right)

	switch expr.Operator.TokenType {
	case BANG:
		return !isTruthy(right)
	case MINUS:
		switch i := right.(type) {
		case float64:
			return -float64(i)
		default:
			fmt.Println("- expected number")
		}
	}

	return nil
}

func (intr Interpreter) VisitVariableExpr(expr Variable) interface{} {
	return intr.Env.get(expr.Name)
}

func (intr Interpreter) VisitAssignExpr(expr Assign) interface{} {
	value := intr.eval(expr.Value)

	intr.Env.assign(expr.Name, value)
	return value
}

func (intr *Interpreter) VisitExpressionStmt(stmt Expression) interface{} {
	intr.eval(stmt.Expression)
	return nil
}

func (intr *Interpreter) VisitPrintStmt(stmt Print) interface{} {
	val := intr.eval(stmt.Expression)
	switch v := val.(type) {
	case string:
		fmt.Println(v)
	default:
		fmt.Println(v)
	}
	return nil
}

func (intr *Interpreter) VisitVarStmt(stmt Var) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = intr.eval(stmt.Initializer)
	}

	intr.Env.define(stmt.Name.Lexeme, value)
	return nil
}

func (intr *Interpreter) VisitBlockStmt(stmt Block) interface{} {
	intr.executeBlock(stmt.Statements, &Environment{Enclosing: intr.Env})
	return nil
}

func (intr *Interpreter) VisitWhileStmt(stmt While) interface{} {
	for isTruthy(intr.eval(stmt.Condition)) {
		intr.execute(stmt.Body)
	}
	return nil
}

func (intr *Interpreter) executeBlock(stmts []Stmt, env *Environment) {
	prev := intr.Env
	defer func() {
		intr.Env = prev
	}()

	intr.Env = env

	for _, stmt := range stmts {
		intr.execute(stmt)
	}
}

func (intr *Interpreter) eval(expr Expr) interface{} {
	v := expr.Accept(intr.visitor())
	//fmt.Println(v)
	return v
}

func (intr Interpreter) execute(stmt Stmt) {
	var v StmtVisitor = &intr
	stmt.Accept(&v)
}

func (intr *Interpreter) Interpret(stmts []Stmt) {
	for _, stmt := range stmts {
		intr.execute(stmt)
	}
}

func (intr *Interpreter) InterpretExpr(expr Expr) {
	val := intr.eval(expr)

	fmt.Println(stringify(val))
}

func (intr *Interpreter) visitor() *ExprVisitor {
	var v ExprVisitor = *intr
	return &v
}

func isTruthy(val interface{}) bool {
	switch i := val.(type) {
	case float64:
		return i != 0
	case string:
		return i != ""
	case bool:
		return bool(i)
	case nil:
		return false
	}
	return true
}

func isEqual(left, right interface{}) bool {
	return left == right
}

func stringify(val interface{}) string {
	return fmt.Sprint(val)
}

// func (intr *Interpreter) print(expr glox.Expr) string {
// 	return fmt.Sprintf("%v", expr.Accept(ast.visitor()))
// }
