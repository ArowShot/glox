package glox

type Expr interface {
    Accept(*ExprVisitor) interface{}
}

type ExprVisitor interface {
    VisitAssignExpr(Assign) interface{}
    VisitBinaryExpr(Binary) interface{}
    VisitGroupingExpr(Grouping) interface{}
    VisitLiteralExpr(Literal) interface{}
    VisitUnaryExpr(Unary) interface{}
    VisitVariableExpr(Variable) interface{}
}

type Assign struct {
    Name Token
    Value Expr
}
func (me Assign) Accept(visitor *ExprVisitor) interface{} {
    v := *visitor
    return v.VisitAssignExpr(me)
}

type Binary struct {
    Left Expr
    Operator Token
    Right Expr
}
func (me Binary) Accept(visitor *ExprVisitor) interface{} {
    v := *visitor
    return v.VisitBinaryExpr(me)
}

type Grouping struct {
    Expression Expr
}
func (me Grouping) Accept(visitor *ExprVisitor) interface{} {
    v := *visitor
    return v.VisitGroupingExpr(me)
}

type Literal struct {
    Value interface{}
}
func (me Literal) Accept(visitor *ExprVisitor) interface{} {
    v := *visitor
    return v.VisitLiteralExpr(me)
}

type Unary struct {
    Operator Token
    Right Expr
}
func (me Unary) Accept(visitor *ExprVisitor) interface{} {
    v := *visitor
    return v.VisitUnaryExpr(me)
}

type Variable struct {
    Name Token
}
func (me Variable) Accept(visitor *ExprVisitor) interface{} {
    v := *visitor
    return v.VisitVariableExpr(me)
}

