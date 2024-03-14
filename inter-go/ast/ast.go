package ast

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
// ‹Literal› ::= ‹number› |  ‹boolean›
// ‹number› ::= ‹digit›+
// ‹digit› ::= 0 |  1 |  2 |  3 |  4 |  5 |  6 |  7 |  8 |  9
// ‹boolean› ::= true |  false
// ‹BinaryExpression› ::= ‹Expression› ‹BinaryOperator› ‹Expression›
// ‹BinaryOperator› ::= + |  - |  * |  / |  == |  <
// ‹UnaryExpression› ::= ‹UnaryOperator› ‹Expression›
// ‹UnaryOperator› ::= ! |  + |  -
// ‹LogicalExpression› ::= ‹Expression› ‹LogicalOperator› ‹Expression›
// ‹LogicalOperator› ::= || |  &&
// ‹ConditionalExpression› ::= ‹Expression› ? ‹Expression› : ‹Expression›

type jsonT map[string]any

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

func newUnaryExpr(expr jsonT) *unaryExpr {
	op := getStr(expr, "operator")
	arg, _ := newExpr(getObj(expr, "argument"))

	return &unaryExpr{op, &arg}
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

func newCondExpr(expr jsonT) *condExpr {
	test, _ := newExpr(getObj(expr, "test"))
	cons, _ := newExpr(getObj(expr, "consequent"))
	altr, _ := newExpr(getObj(expr, "alternate"))

	return &condExpr{&test, &cons, &altr}
}

func newBindExpr(bind jsonT) *bindExpr {
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

	return &bindExpr{binds}
}

func newBlockExpr(body []any) *blockExpr {
	if len(body) < 1 {
		// this is shite
		return nil
	}

	exprs := make([]expr, len(body))
	var nextExpr expr
	var err error
	for i, expr := range body {
		nextExpr, err = newExpr(expr.(map[string]any))
		if err != nil {
			panic(err.Error())
		}
		exprs[i] = nextExpr

	}

	return &blockExpr{exprs}

	// first := body[0].(map[string]any)
	//first := body[0]
	// var resExpr expr
	// switch getStr(first, "type") {
	// case "VariableDeclaration":
	// 	resExpr = newBindExpr(first, body[1:])
	//
	// case "ExpressionStatement":
	// 	fallthrough
	// case "ReturnStatement":
	// 	resExpr, _ = newExpr(first)
	// }
	//
	// return resExpr
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

func newCallExpr(expr jsonT) *callExpr {
	callee, _ := newExpr(getObj(expr, "callee"))
	arg, _ := newExpr(getArr(expr, "arguments")[0].(map[string]any))

	return &callExpr{&callee, &arg}
}

func newLitExpr(expr jsonT) *litExpr {
	var (
		raw     = getStr(expr, "raw")
		val     any
		valType string
	)

	if strings.ContainsAny(raw, "0123456789") {
		val, _ = strconv.ParseInt(raw, 10, 64)
		valType = number
	} else {
		val, _ = strconv.ParseBool(raw)
		valType = boolean
	}

	return &litExpr{val, valType}
}

func newAssignExpr(expr jsonT) *assignExpr {
	op := getStr(expr, "operator")
	lhs, _ := newExpr(getObj(expr, "left"))
	ident := lhs.(*refExpr).string
	rhs, _ := newExpr(getObj(expr, "right"))
	return &assignExpr{op, ident, rhs}
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
	case "VariableDeclaration":
		resExpr = newBindExpr(json)
	case "BlockStatement":
		resExpr = newBlockExpr(getArr(json, "body"))
	case "AssignmentExpression":
		resExpr = newAssignExpr(json)
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

type program struct {
	body []expr
}

func New(json jsonT) (*program, error) {
	body := json["body"].([]any)
	exprs := make([]expr, 0, 5)
	for _, expr := range body {
		newExpr, err := newExpr(expr.(map[string]any))
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, newExpr)
	}

	// prog := newBlockExpr(body)
	return &program{exprs}, nil
}
