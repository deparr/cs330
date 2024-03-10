package ast

import "fmt"

// really should rethink type aei
type Value struct {
	self any
}

type fnValue struct {
	env  environ
	arg  string
	body *expr
}

type void struct{}

func (val *Value) getType() string {
	switch (*val).self.(type) {
	case bool:
		return BOOLEAN
	case int64:
		return NUMBER
	case fnValue:
		return FUNCTION
	case void:
		return VOID
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

func (bx binExpr) Eval(env environ, heap []Value) (*Value, error) {
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

func (ux *unaryExpr) Eval(env environ, heap []Value) (*Value, error) {
	arg, err := (*ux.arg).Eval(env, heap)
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

func (lx *logicalExpr) Eval(env environ, heap []Value) (*Value, error) {
	left, err := (*lx.lhs).Eval(env, heap)
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

	right, err = (*lx.rhs).Eval(env, heap)
	if err != nil {
		return nil, err
	}
	if right.getType() != BOOLEAN {
		return nil, fmt.Errorf("Got non-bool operand in LogOp: `%s`", lx.op)
	}

	rightBool = right.getBool()
	return &Value{self: rightBool}, nil
}

func (cx *condExpr) Eval(env environ, heap []Value) (*Value, error) {
	test, err := (*cx.test).Eval(env, heap)
	if err != nil {
		return nil, err
	}

	if test.getType() != BOOLEAN {
		return nil, fmt.Errorf("Non bool in conditonal")
	}

	if test.getBool() {
		return (*cx.cons).Eval(env, heap)
	}

	return (*cx.altr).Eval(env, heap)
}

func (bx *bindExpr) Eval(env environ, heap []Value) (*Value, error) {
	for _, bind := range bx.binds {
		boundValue, err := bind.expr.Eval(env, heap)
		if err != nil {
			return nil, err
		}

		env = env.extend(bind.string, *boundValue)
	}

	return (*bx.body).Eval(env, heap)
}

func (rx *refExpr) Eval(env environ, heap []Value) (*Value, error) {
	val, bound := env.lookup(rx.string)
	if !bound {
		return nil, fmt.Errorf("Unbound identifier: `%s`", rx.string)
	}

	return &val, nil
}

func (fx fnExpr) Eval(env environ, heap []Value) (*Value, error) {
	return &Value{fnValue{env, fx.arg, fx.body}}, nil
}

func (cx callExpr) Eval(env environ, heap []Value) (*Value, error) {
	fnVal, err := (*cx.callee).Eval(env, heap)
	if err != nil {
		return nil, err
	}

	if fnVal.getType() != FUNCTION {
		return nil, fmt.Errorf("`%s` is not callable", fnVal.getType())
	}
	// I really just wish I could shadow vars in go
	fn := fnVal.self.(fnValue)
	argVal, err := (*cx.arg).Eval(env, heap)
	if err != nil {
		return nil, err
	}

	return (*fn.body).Eval(env.extend(fn.arg, *argVal), heap)
}

func (cx *litExpr) Eval(env environ, heap []Value) (*Value, error) {
	return &Value{self: cx.val}, nil
}
