package ast

// put type defs and strs ghere
import (
	"fmt"
	"strings"
)

const (
	// TYPES
	NUMBER   = "number"
	BOOLEAN  = "boolean"
	FUNCTION = "function"
	VOID     = "void"

	// OPERATORS
	PLUS  = "+"
	MINUS = "-"
	MUL   = "*"
	DIV   = "/"
	AND   = "&&"
	OR    = "||"
	EQL   = "=="
	NEQL  = "!="
	LT    = "<"
	GT    = ">"
	LTE   = "<="
	GTE   = ">="
	NOT   = "!"
)

type expr interface {
	Eval(env environ, heap []Value) (*Value, error)
}

// / Binary Expressions -------------------------------------
type binExpr struct {
	op  string // todo: not sure about this
	lhs *expr
	rhs *expr
}

// / Unary Expressions -------------------------------------
type unaryExpr struct {
	op  string
	arg *expr
}

// / Logical Expressions ------------------------------------
type logicalExpr struct {
	op  string
	lhs *expr
	rhs *expr
}

// / Conditional Expressions
type condExpr struct {
	test *expr
	cons *expr
	altr *expr
}

type bindExpr struct {
	binds []struct {
		string
		expr
	}
	body *expr
}

type refExpr struct {
	string
}

// / Function Expressions -----------------------------------
type fnExpr struct {
	arg  string
	body *expr
}

// / Call Expressions ---------------------------------------
type callExpr struct {
	callee *expr
	arg    *expr
}

// / Literal Expressions ------------------------------------
type litExpr struct {
	val     any
	valType string
}

type assignExpr struct {
	op    string
	left  expr
	right expr
}

func (bx binExpr) String() string {
	var _type string
	if bx.op == PLUS || bx.op == MINUS || bx.op == DIV || bx.op == MUL {
		_type = "arithmetic"
	} else {
		_type = "relational"
	}
	return fmt.Sprintf("(%s %s %s %s)", _type, bx.op, *bx.lhs, *bx.rhs)
}

func (ux unaryExpr) String() string {
	return fmt.Sprintf("(unary %s %s)", ux.op, *ux.arg)
}

func (lx logicalExpr) String() string {
	return fmt.Sprintf("(logical %s %s %s)", lx.op, *lx.lhs, *lx.rhs)
}

func (cx condExpr) String() string {
	return fmt.Sprintf("(conditional %s %s %s)", *cx.test, *cx.cons, *cx.altr)
}

func (bx bindExpr) String() string {
	bindStrs := make([]string, len(bx.binds))
	for i, bind := range bx.binds {
		bindStrs[i] = fmt.Sprintf("[%s %s]", bind.string, bind.expr)
	}
	return fmt.Sprintf("(let %s %s)", strings.Join(bindStrs, " "), *bx.body)
}

func (rx refExpr) String() string {
	//return fmt.Sprintf("(ref %s)", rx.string)
	return rx.string
}

func (fx fnExpr) String() string {
	return fmt.Sprintf("(%s) %s", fx.arg, *fx.body)
}

func (cx callExpr) String() string {
	return fmt.Sprintf("(%s %s)", *cx.callee, *cx.arg)
}

func (lx litExpr) String() string {
	if lx.valType == NUMBER {
		return fmt.Sprintf("(%s %d)", lx.valType, lx.val.(int64))
	}
	return fmt.Sprintf("(%s %t)", lx.valType, lx.val.(bool))
}
