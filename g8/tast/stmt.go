package tast

// ExprStmt is a statement with just an expression.
type ExprStmt struct {
	Expr
}

// AssignStmt is an assignment statement, like "a,b=x,y".
type AssignStmt struct {
	Left  *ExprList
	Right *ExprList
}
