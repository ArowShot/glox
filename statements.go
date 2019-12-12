package glox

type Stmt interface {
    Accept(*StmtVisitor) interface{}
}

type StmtVisitor interface {
    VisitBlockStmt(Block) interface{}
    VisitExpressionStmt(Expression) interface{}
    VisitPrintStmt(Print) interface{}
    VisitVarStmt(Var) interface{}
    VisitWhileStmt(While) interface{}
}

type Block struct {
    Statements []Stmt
}
func (me Block) Accept(visitor *StmtVisitor) interface{} {
    v := *visitor
    return v.VisitBlockStmt(me)
}

type Expression struct {
    Expression Expr
}
func (me Expression) Accept(visitor *StmtVisitor) interface{} {
    v := *visitor
    return v.VisitExpressionStmt(me)
}

type Print struct {
    Expression Expr
}
func (me Print) Accept(visitor *StmtVisitor) interface{} {
    v := *visitor
    return v.VisitPrintStmt(me)
}

type Var struct {
    Name Token
    Initializer Expr
}
func (me Var) Accept(visitor *StmtVisitor) interface{} {
    v := *visitor
    return v.VisitVarStmt(me)
}

type While struct {
    Condition Expr
    Body Stmt
}
func (me While) Accept(visitor *StmtVisitor) interface{} {
    v := *visitor
    return v.VisitWhileStmt(me)
}

