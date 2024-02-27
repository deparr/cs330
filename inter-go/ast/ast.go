package ast

// TODO: split this into multiple files, rip

import (
	"fmt"
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
	NUMBER   = "number"
	BOOLEAN  = "boolean"
	FUNCTION = "function"

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

type fnValue struct {
	env  environ
	arg  string
	body *expr
}

func (fn fnValue) String() string {
	return "(function)"
}

type environ struct {
	env map[string]Value
}

func EmptyEnv() environ {
	return environ{env: make(map[string]Value)}
}

func (env environ) extend(newKey string, newValue Value) environ {
	new := EmptyEnv()
	for k, v := range env.env {
		new.env[k] = v
	}
	new.env[newKey] = newValue
	return new
}

func (env environ) lookup(ident string) (Value, bool) {
	val, bound := env.env[ident]
	return val, bound
}

func (val *Value) getType() string {
	switch (*val).self.(type) {
	case bool:
		return BOOLEAN
	case int64:
		return NUMBER
	case fnValue:
		return FUNCTION
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
	_type := v.getType()
	if _type == NUMBER {
		return fmt.Sprintf("(number %d)", v.getNumber())
	} else if _type == BOOLEAN {
		return fmt.Sprintf("(boolean %t)", v.getBool())
	} else if _type == FUNCTION {
		return "(function)"
	}
	return _type
}

type expr interface {
	Eval(env environ) (*Value, error)
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

func (bx binExpr) Eval(env environ) (*Value, error) {
	left, err := (*bx.lhs).Eval(env)
	if err != nil {
		return nil, err
	}
	right, err := (*bx.rhs).Eval(env)
	if err != nil {
		return nil, err
	}

	leftT := left.getType()
	rightT := right.getType()

	if leftT != NUMBER || rightT != NUMBER {
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

func (ux *unaryExpr) Eval(env environ) (*Value, error) {
	arg, err := (*ux.arg).Eval(env)
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

func (lx *logicalExpr) Eval(env environ) (*Value, error) {
	left, err := (*lx.lhs).Eval(env)
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

	right, err = (*lx.rhs).Eval(env)
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

func (cx *condExpr) Eval(env environ) (*Value, error) {
	test, err := (*cx.test).Eval(env)
	if err != nil {
		return nil, err
	}

	if test.getType() != BOOLEAN {
		return nil, fmt.Errorf("Non bool in conditonal")
	}

	if test.getBool() {
		return (*cx.cons).Eval(env)
	}

	return (*cx.altr).Eval(env)
}

func (cx condExpr) String() string {
	return fmt.Sprintf("(conditional %s %s %s)", *cx.test, *cx.cons, *cx.altr)
}

type bindExpr struct {
	binds []struct {
		string
		expr
	}
	body *expr
}

func newBindExpr(bind jsonT, rest []any) *bindExpr {
	binds := make([]struct {
		string
		expr
	}, 0, 3)

	for _, dec := range getArr(bind, "declarations") {
		decObj := dec.(map[string]any)
		ident := getStr(getObj(decObj, "id"), "name")
		init, _ := newExpr(getObj(decObj, "init"))
		binds = append(binds, struct {
			string
			expr
		}{ident, init})
	}

	if len(rest) < 1 {
		panic("newBindExpr()::rest has zero len")
	}

	bodyObj := rest[0].(map[string]any)
	var body expr
	if getStr(bodyObj, "type") == "VariableDeclaration" {
		body = newBindExpr(bodyObj, rest[1:])
	} else {
		body, _ = newExpr(bodyObj)
	}

	return &bindExpr{binds, &body}
}

func (bx bindExpr) String() string {
	bindStrs := make([]string, len(bx.binds))
	for i, bind := range bx.binds {
		bindStrs[i] = fmt.Sprintf("[%s %s]", bind.string, bind.expr)
	}
	return fmt.Sprintf("(let %s %s)", strings.Join(bindStrs, " "), *bx.body)
}

func (bx *bindExpr) Eval(env environ) (*Value, error) {
	for _, bind := range bx.binds {
		boundValue, err := bind.expr.Eval(env)
		if err != nil {
			return nil, err
		}

		env = env.extend(bind.string, *boundValue)
	}

	return (*bx.body).Eval(env)
}

func newBlockExpr(body []any) expr {
	if len(body) < 1 {
		return nil
	}

	first := body[0].(map[string]any)
	//first := body[0]
	var resExpr expr
	switch getStr(first, "type") {
	case "VariableDeclaration":
		resExpr = newBindExpr(first, body[1:])

	case "ExpressionStatement":
		fallthrough
	case "ReturnStatement":
		resExpr, _ = newExpr(first)
	}

	return resExpr
}

type refExpr struct {
	string
}

func (rx *refExpr) Eval(env environ) (*Value, error) {
	val, bound := env.lookup(rx.string)
	if !bound {
		return nil, fmt.Errorf("Unbound identifier: `%s`", rx.string)
	}

	return &val, nil
}

func (rx refExpr) String() string {
	//return fmt.Sprintf("(ref %s)", rx.string)
	return rx.string
}

// / Function Expressions -----------------------------------
type fnExpr struct {
	arg  string
	body *expr
}

func newFnExpr(expr jsonT) *fnExpr {
	params := getArr(expr, "params")
	var arg string = ""
	if len(params) > 0 {
		arg = getStr((params[0].(map[string]any)), "name")
	}
	body, _ := newExpr(getObj(expr, "body"))

	return &fnExpr{arg, &body}
}

func (fx fnExpr) Eval(env environ) (*Value, error) {
	return &Value{fnValue{env, fx.arg, fx.body}}, nil
}

func (fx fnExpr) String() string {
	return fmt.Sprintf("(%s) %s", fx.arg, *fx.body)
}

// / Call Expressions ---------------------------------------
type callExpr struct {
	callee *expr
	arg    *expr
}

func newCallExpr(expr jsonT) *callExpr {
	callee, _ := newExpr(getObj(expr, "callee"))
	arg, _ := newExpr(getArr(expr, "arguments")[0].(map[string]any))

	return &callExpr{&callee, &arg}
}

func (cx callExpr) String() string {
	return fmt.Sprintf("(%s %s)", *cx.callee, *cx.arg)
}

func (cx callExpr) Eval(env environ) (*Value, error) {
	fnVal, err := (*cx.callee).Eval(env)
	if err != nil {
		return nil, err
	}

	if fnVal.getType() != FUNCTION {
		return nil, fmt.Errorf("`%s` is not callable", fnVal.getType())
	}
	// I really just wish I could shadow vars in go
	fn := fnVal.self.(fnValue)
	argVal, err := (*cx.arg).Eval(env)
	if err != nil {
		return nil, err
	}

	return (*fn.body).Eval(env.extend(fn.arg, *argVal))
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

func (cx *litExpr) Eval(env environ) (*Value, error) {
	return &Value{self: cx.val}, nil
}

func (lx litExpr) String() string {
	if lx.valType == NUMBER {
		return fmt.Sprintf("(%s %d)", lx.valType, lx.val.(int64))
	}
	return fmt.Sprintf("(%s %t)", lx.valType, lx.val.(bool))
}

// / All Expressions ----------------------------------------
// TODO: dont ignore errors
func newExpr(json jsonT) (expr, error) {
	exprType := json["type"].(string)
	var resExpr expr
	switch exprType {
	case "UnaryExpression":
		resExpr = newUnaryExpr(json)
	case "BinaryExpression":
		resExpr = newBinExpr(json)
	case "LogicalExpression":
		resExpr = newLogicalExpr(json)
	case "ConditionalExpression":
		resExpr = newCondExpr(json)
	case "FunctionExpression":
		resExpr = newFnExpr(json)
	case "CallExpression":
		resExpr = newCallExpr(json)
	case "BlockStatement":
		resExpr = newBlockExpr(getArr(json, "body"))
	case "Literal":
		resExpr = newLitExpr(json)
	case "Identifier":
		resExpr = &refExpr{getStr(json, "name")}
	case "ExpressionStatement":
		resExpr, _ = newExpr(getObj(json, "expression"))
	case "ReturnStatement":
		resExpr, _ = newExpr(getObj(json, "argument"))
	default:
		return nil, fmt.Errorf("Unknown type in newExpr(): %s", exprType)
	}

	return resExpr, nil
}

func New(json jsonT) (expr, error) {
	body := json["body"].([]any)

	prog := newBlockExpr(body)
	return prog, nil
}
