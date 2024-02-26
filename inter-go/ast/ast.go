package ast

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ‹Program› ::= ‹Statement›
// ‹Statement› ::= ‹ExpressionStatement›
// ‹ExpressionStatement› ::= ‹Expression›
// ‹Expression› ::= ‹Literal›
// 			  |  ‹BinaryExpression›
// 			  |  ‹UnaryExpression›
// 			  |  ‹LogicalExpression›
// 			  |  ‹ConditionalExpression›
// ‹Literal› ::= ‹number›
// 			  |  ‹boolean›
// ‹number› ::= ‹digit›+
// ‹digit› ::= 0
// 			|  1
// 			|  2
// 			|  3
// 			|  4
// 			|  5
// 			|  6
// 			|  7
// 			|  8
// 			|  9
// ‹boolean› ::= true |  false
// ‹BinaryExpression› ::= ‹Expression› ‹BinaryOperator› ‹Expression›
// ‹BinaryOperator› ::= + |  - |  * |  / |  == |  <
// ‹UnaryExpression› ::= ‹UnaryOperator› ‹Expression›
// ‹UnaryOperator› ::= ! |  + |  -
// ‹LogicalExpression› ::= ‹Expression› ‹LogicalOperator› ‹Expression›
// ‹LogicalOperator› ::= || |  &&
// ‹ConditionalExpression› ::= ‹Expression› ? ‹Expression› : ‹Expression›

type jsonT map[string]any

const (
	// TYPES
	NUMBER  = "number"
	BOOLEAN = "boolean"

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

func getArr(json jsonT, key string) []any {
	return json[key].([]any)
}

func getObj(json jsonT, key string) map[string]any {
	return json[key].(map[string]any)
}

func getStr(json jsonT, key string) string {
	return json[key].(string)
}

func getNum(json jsonT, key string) int64 {
	return json[key].(int64)
}

type Value struct {
	self any
}

func (val *Value) getType() string {
	switch (*val).self.(type) {
	case bool:
		return BOOLEAN
	case int64:
		return NUMBER
	default:
		return "unknown"
	}
}

func (val *Value) getNumber() int64 {
	return (*val).self.(int64)
}

func (val *Value) getBool() bool {
	return (*val).self.(bool)
}

func (v Value) String() string {
	if v.getType() == NUMBER {
		return fmt.Sprintf("(number %d)", v.getNumber())
	}
	return fmt.Sprintf("(boolean %t)", v.getBool())
}

type expr interface {
	Eval() (*Value, error)
}

// / Binary Expressions -------------------------------------
type binExpr struct {
	op  string // todo: not sure about this
	lhs *expr
	rhs *expr
}

func newBinExpr(expr jsonT) *binExpr {
	op := getStr(expr, "operator")
	lhs, _ := newExpr(getObj(expr, "left"))
	rhs, _ := newExpr(getObj(expr, "right"))

	return &binExpr{
		op,
		&lhs,
		&rhs,
	}
}

func (bx binExpr) Eval() (*Value, error) {
	left, err := (*bx.lhs).Eval()
	if err != nil {
		return nil, err
	}
	right, err := (*bx.rhs).Eval()
	if err != nil {
		return nil, err
	}

	leftT := left.getType()
	rightT := right.getType()

	if leftT != NUMBER || rightT != NUMBER {
		fmt.Fprintf(os.Stderr, "%s %s %s\n", leftT, rightT, bx.op)
		return nil, fmt.Errorf("Got non-number operand in BinOp `%s`", bx.op)
	}

	leftNum := left.getNumber()
	rightNum := right.getNumber()

	var result any
	switch bx.op {
	case PLUS:
		result = leftNum + rightNum
	case MINUS:
		result = leftNum - rightNum
	case MUL:
		result = leftNum * rightNum
	case DIV:
		if rightNum == 0 {
			return nil, fmt.Errorf("Division by zero")
		}
		result = leftNum / rightNum
	case EQL:
		result = leftNum == rightNum
	case NEQL:
		result = leftNum != rightNum
	case LT:
		result = leftNum < rightNum
	case LTE:
		result = leftNum <= rightNum
	case GT:
		result = leftNum > rightNum
	case GTE:
		result = leftNum >= rightNum
	default:
		return nil, fmt.Errorf("Unknown BinOp: %s", bx.op)
	}

	return &Value{self: result}, nil
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

// / Unary Expressions -------------------------------------
type unaryExpr struct {
	op  string
	arg *expr
}

func newUnaryExpr(expr jsonT) *unaryExpr {
	op := getStr(expr, "operator")
	arg, _ := newExpr(getObj(expr, "argument"))

	return &unaryExpr{op, &arg}
}

func (ux *unaryExpr) Eval() (*Value, error) {
	arg, err := (*ux.arg).Eval()
	if err != nil {
		return nil, err
	}

	argT := arg.getType()
	var result any
	switch ux.op {
	case NOT:
		if argT != BOOLEAN {
			return nil, fmt.Errorf("Expected bool operand with `!` operator, got: %s", argT)
		}
		result = !arg.getBool()

	case PLUS:
		if argT != NUMBER {
			return nil, fmt.Errorf("Expected number operand with unary `+` oeprator, got %s", argT)
		}
		if arg.getNumber() < 0 {
			result = arg.getNumber() * -1
		} else {
			result = arg.getNumber()
		}

	case MINUS:
		if argT != NUMBER {
			return nil, fmt.Errorf("Expected number operand with unary `-` oeprator, got %s", argT)
		}
		result = arg.getNumber() * 1
	}

	return &Value{self: result}, nil
}

func (ux unaryExpr) String() string {
	return fmt.Sprintf("(unary %s %s)", ux.op, *ux.arg)
}

// / Logical Expressions ------------------------------------
type logicalExpr struct {
	op  string
	lhs *expr
	rhs *expr
}

func newLogicalExpr(expr jsonT) *logicalExpr {
	op := getStr(expr, "operator")
	lhs, _ := newExpr(getObj(expr, "left"))
	rhs, _ := newExpr(getObj(expr, "right"))

	return &logicalExpr{
		op,
		&lhs,
		&rhs,
	}
}

func (lx *logicalExpr) Eval() (*Value, error) {
	left, err := (*lx.lhs).Eval()
	if err != nil {
		return nil, err
	}

	if left.getType() != BOOLEAN {
		return nil, fmt.Errorf("Got non-bool operand in LogOp `%s`", lx.op)
	}

	var (
		leftBool  = left.getBool()
		right     *Value
		rightBool bool
	)
	// Short circuits
	switch lx.op {
	case OR:
		if leftBool {
			return &Value{self: true}, nil
		}
	case AND:
		if !leftBool {
			return &Value{self: false}, nil
		}
	default:
		return nil, fmt.Errorf("Unimplemented logical op: `%s`", lx.op)
	}

	right, err = (*lx.rhs).Eval()
	if err != nil {
		return nil, err
	}
	if right.getType() != BOOLEAN {
		return nil, fmt.Errorf("Got non-bool operand in LogOp: `%s`", lx.op)
	}

	rightBool = right.getBool()
	return &Value{self: rightBool}, nil
}

func (lx logicalExpr) String() string {
	return fmt.Sprintf("(logical %s %s %s)", lx.op, *lx.lhs, *lx.rhs)
}

// / Conditional Expressions
type condExpr struct {
	test *expr
	cons *expr
	altr *expr
}

func newCondExpr(expr jsonT) *condExpr {
	test, _ := newExpr(getObj(expr, "test"))
	cons, _ := newExpr(getObj(expr, "consequent"))
	altr, _ := newExpr(getObj(expr, "alternate"))

	return &condExpr{&test, &cons, &altr}
}

func (cx *condExpr) Eval() (*Value, error) {
	test, err := (*cx.test).Eval()
	if err != nil {
		return nil, err
	}

	if test.getType() != BOOLEAN {
		return nil, fmt.Errorf("Non bool in conditonal")
	}

	if test.getBool() {
		return (*cx.cons).Eval()
	}

	return (*cx.altr).Eval()
}

func (cx condExpr) String() string {
	return fmt.Sprintf("(conditional %s %s %s)", *cx.test, *cx.cons, *cx.altr)
}

// / Literal Expressions ------------------------------------
type litExpr struct {
	val     any
	valType string
}

func newLitExpr(expr jsonT) *litExpr {
	var (
		raw     = getStr(expr, "raw")
		val     any
		valType string
	)

	if strings.ContainsAny(raw, "0123456789") {
		val, _ = strconv.ParseInt(raw, 10, 64)
		valType = NUMBER
	} else {
		val, _ = strconv.ParseBool(raw)
		valType = BOOLEAN
	}

	return &litExpr{val, valType}
}

func (cx *litExpr) Eval() (*Value, error) {
	return &Value{self: cx.val}, nil
}

func (lx litExpr) String() string {
	if lx.valType == NUMBER {
		return fmt.Sprintf("(%s %d)", lx.valType, lx.val.(int64))
	}
	return fmt.Sprintf("(%s %t)", lx.valType, lx.val.(bool))
}

// / All Expressions ----------------------------------------
func newExpr(json jsonT) (expr, error) {
	exprType := json["type"].(string)
	var newExpr expr
	switch exprType {
	case "UnaryExpression":
		newExpr = newUnaryExpr(json)
	case "BinaryExpression":
		newExpr = newBinExpr(json)
	case "LogicalExpression":
		newExpr = newLogicalExpr(json)
	case "ConditionalExpression":
		newExpr = newCondExpr(json)
	case "Literal":
		newExpr = newLitExpr(json)
	default:
		return nil, fmt.Errorf("Unknown type in newExpr(): %s", exprType)
	}

	return newExpr, nil
}

func New(json jsonT) (expr, error) {
	body := json["body"].([]any)
	// i := 0
	// for i < len(body) {
	// }

	first := body[0].(map[string]any)
	first = getObj(first, "expression")
	return newExpr(first)
}
