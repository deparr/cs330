package ast

import "fmt"

// TODO: this type is bad, should be an interface probably
type Value struct {
	self any
}

type fnValue struct {
	env  environ
	arg  string
	body *expr
}

type voidValue struct{}

func (val *Value) getType() string {
	switch (*val).self.(type) {
	case bool:
		return boolean
	case int64:
		return number
	case fnValue:
		return funciton
	case voidValue:
		return void
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
	if _type == number {
		return fmt.Sprintf("(number %d)", v.getNumber())
	} else if _type == boolean {
		return fmt.Sprintf("(boolean %t)", v.getBool())
	} else if _type == funciton {
		return "(function)"
	}
	return _type
}

func (fn fnValue) String() string {
	return "(function)"
}

type environ struct {
	env map[string]int
}

func EmptyEnv() environ {
	return environ{env: make(map[string]int)}
}

func (env environ) extend(newKey string, newValue Value, heap *[]Value) environ {
	new := EmptyEnv()
	for k, v := range env.env {
		new.env[k] = v
	}
	*heap = append(*heap, newValue)
	new.env[newKey] = len(*heap) - 1
	return new
}

func (env environ) lookup(ident string, heap *[]Value) (Value, bool, int) {
	addr, bound := env.env[ident]
	var val Value
	if addr < len(*heap) {
		val = (*heap)[addr]
	}
	return val, bound, addr
}

func (prog program) Eval() (*Value, error) {
	env, heap := EmptyEnv(), []Value{}
	var ret *Value
	var err error
	for _, stat := range prog.body {
		fmt.Printf("Evaling: %s with env/heap %v/%s\n", stat, env.env, heap)
		ret, err = stat.Eval(&env, &heap)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (bx binExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	left, err := (*bx.lhs).Eval(env, heap)
	if err != nil {
		return nil, err
	}
	right, err := (*bx.rhs).Eval(env, heap)
	if err != nil {
		return nil, err
	}

	leftT := left.getType()
	rightT := right.getType()

	if leftT != number || rightT != number {
		return nil, fmt.Errorf("Got non-number operand in BinOp `%s`", bx.op)
	}

	leftNum := left.getNumber()
	rightNum := right.getNumber()

	var result any
	switch bx.op {
	case plus:
		result = leftNum + rightNum
	case minus:
		result = leftNum - rightNum
	case mul:
		result = leftNum * rightNum
	case div:
		if rightNum == 0 {
			return nil, fmt.Errorf("Division by zero")
		}
		result = leftNum / rightNum
	case eql:
		result = leftNum == rightNum
	case neql:
		result = leftNum != rightNum
	case lt:
		result = leftNum < rightNum
	case lte:
		result = leftNum <= rightNum
	case gt:
		result = leftNum > rightNum
	case gte:
		result = leftNum >= rightNum
	default:
		return nil, fmt.Errorf("Unknown BinOp: %s", bx.op)
	}

	return &Value{self: result}, nil
}

func (ux *unaryExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	arg, err := (*ux.arg).Eval(env, heap)
	if err != nil {
		return nil, err
	}

	argT := arg.getType()
	var result any
	switch ux.op {
	case not:
		if argT != boolean {
			return nil, fmt.Errorf("Expected bool operand with `!` operator, got: %s", argT)
		}
		result = !arg.getBool()

	case plus:
		if argT != number {
			return nil, fmt.Errorf("Expected number operand with unary `+` oeprator, got %s", argT)
		}
		if arg.getNumber() < 0 {
			result = arg.getNumber() * -1
		} else {
			result = arg.getNumber()
		}

	case minus:
		if argT != number {
			return nil, fmt.Errorf("Expected number operand with unary `-` oeprator, got %s", argT)
		}
		result = arg.getNumber() * 1
	}

	return &Value{self: result}, nil
}

func (lx *logicalExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	left, err := (*lx.lhs).Eval(env, heap)
	if err != nil {
		return nil, err
	}

	if left.getType() != boolean {
		return nil, fmt.Errorf("Got non-bool operand in LogOp `%s`", lx.op)
	}

	var (
		leftBool  = left.getBool()
		right     *Value
		rightBool bool
	)
	// Short circuits
	switch lx.op {
	case or:
		if leftBool {
			return &Value{self: true}, nil
		}
	case and:
		if !leftBool {
			return &Value{self: false}, nil
		}
	default:
		return nil, fmt.Errorf("Unimplemented logical op: `%s`", lx.op)
	}

	right, err = (*lx.rhs).Eval(env, heap)
	if err != nil {
		return nil, err
	}
	if right.getType() != boolean {
		return nil, fmt.Errorf("Got non-bool operand in LogOp: `%s`", lx.op)
	}

	rightBool = right.getBool()
	return &Value{self: rightBool}, nil
}

func (cx *condExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	test, err := (*cx.test).Eval(env, heap)
	if err != nil {
		return nil, err
	}

	if test.getType() != boolean {
		return nil, fmt.Errorf("Non bool in conditonal")
	}

	if test.getBool() {
		return (*cx.cons).Eval(env, heap)
	}

	return (*cx.altr).Eval(env, heap)
}

func (bx *bindExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	for _, bind := range bx.binds {
		boundValue, err := bind.expr.Eval(env, heap)
		if err != nil {
			return nil, err
		}

		*env = env.extend(bind.string, *boundValue, heap)
	}

	return &Value{voidValue{}}, nil
}

func (rx refExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	val, bound, _ := env.lookup(rx.string, heap)
	if !bound {
		return nil, fmt.Errorf("Unbound identifier: `%s`", rx.string)
	}

	return &val, nil
}

func (fx fnExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	return &Value{fnValue{*env, fx.arg, fx.body}}, nil
}

func (cx callExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	fnVal, err := (*cx.callee).Eval(env, heap)
	if err != nil {
		return nil, err
	}

	if fnVal.getType() != funciton {
		return nil, fmt.Errorf("`%s` is not callable", fnVal.getType())
	}
	// I really just wish I could shadow vars in go
	fn := fnVal.self.(fnValue)
	argVal, err := (*cx.arg).Eval(env, heap)
	if err != nil {
		return nil, err
	}

	newEnv := env.extend(fn.arg, *argVal, heap)
	return (*fn.body).Eval(&newEnv, heap)
}

func (cx *litExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	return &Value{self: cx.val}, nil
}

func (ax *assignExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	switch ax.op {
	case assign:
		_, bound, addr := env.lookup(ax.ident, heap)
		if !bound {
			return nil, fmt.Errorf("Attempt to assign unbound identifer: `%s`", ax.ident)
		}

		rhs, err := ax.rhs.Eval(env, heap)
		if err != nil {
			return nil, err
		}

		(*heap)[addr] = *rhs
	}
	return &Value{self: voidValue{}}, nil
}

func (bx blockExpr) Eval(env *environ, heap *[]Value) (*Value, error) {
	var res *Value
	var err error
	for _, expr := range bx.exprs {
		res, err = expr.Eval(env, heap)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
