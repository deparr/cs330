package types

import (
	"fmt"
	"strings"
)

// I hate all of this
// this code is shite
// this is some of the worst code I have ever written

type environ map[string]ast_t

func Parse(json map[string]any) (ast_t, error) {
	return parse(json, environ{})
}

func parse(json map[string]any, env environ) (ast_t, error) {
	item_t := json["type"].(string)
	if item_t != "list" && item_t != "symbol" {
		return newLiteral(item_t)
	} else if item_t == "symbol" {
		ident := json["value"].(string)
		_type, bound := env[ident]
		if !bound {
			return nil, fmt.Errorf("Cannot access undefined ident %s", ident)
		}
		return _type, nil
	}

	var (
		items  = json["value"].([]any)
		i      = 0
		item_v string
		ast    ast_t
		err    error
	)

	for i < len(items) {
		item := items[i].(map[string]any)
		item_t = item["type"].(string)
		// fmt.Printf("%d %v\n", i, item)
		// if type is not list ??
		// what to do here?
		// can't switch on value because that can be any type
		// WARN: this probably shouldn't be here
		item_v = item["value"].(string)
		switch item_v {
		case "object":
			ast, err = parseObjLit(items[i+1:], env)
			if err != nil {
				return nil, err
			}
			i++
			obj_len := len((ast.(ast_obj)).Fields)
			i += obj_len

		case "field":
			obj, err := parse(items[i+1].(map[string]any), env)
			if err != nil {
				return nil, err
			}

			if !strings.Contains(obj.Type(), "obj{") {
				return nil, fmt.Errorf("'%s' type does not have field access", obj.Type())
			}

			ident := items[i+2].(map[string]any)
			ident_t := ident["type"].(string)
			if ident_t != "symbol" {
				return nil, fmt.Errorf("expected object field of type symbol, got: %s", ident_t)
			}

			ident_str := ident["value"].(string)
			ast_as_obj := obj.(ast_obj)
			field_type, prs := ast_as_obj.Fields[ident_str]
			if !prs {
				return nil, fmt.Errorf("type `%s` has no field `%s`", ast_as_obj, ident_str)
			}

			return field_type, nil

		case "fun":
			arg_l := (items[i+1].(map[string]any))["value"].([]any)
			ident_str := (arg_l[0].(map[string]any))["value"].(string)
			arg_type, err := newLiteral(arg_l[2].(map[string]any)["value"].(string))
			if err != nil {
				return nil, err
			}

			old_type, was_bound := env[ident_str]
			env[ident_str] = arg_type
			body_type, err := parse(items[i+2].(map[string]any), env)
			if err != nil {
				return nil, err
			}

			if was_bound {
				env[ident_str] = old_type
			} else {
				delete(env, ident_str)
			}

			return ast_func{arg_type, body_type}, nil

		case "app":
			fun, err := parse(items[i+1].(map[string]any), env)
			if err != nil {
				return nil, err

			}

			if !strings.Contains(fun.Type(), "func|") {
				return nil, fmt.Errorf("`%s` type is not callable", fun.Type())
			}

			arg, err := parse(items[i+2].(map[string]any), env)
			if err != nil {
				return nil, err
			}
			fun_as_func := fun.(ast_func)

			if arg.Type() != fun_as_func.Arg.Type() {
				return nil, fmt.Errorf("function `%s` cannot accept arg of type: %s", fun, arg)
			}

			return fun_as_func.Ret, nil

		case "let":
			bind_l := (items[i+1].(map[string]any))["value"].([]any)
			binds := make([]struct {
				string
				ast_t
				bool
			}, len(bind_l))
			for i, bind := range bind_l {
				bind_arg_l := bind.(map[string]any)
				bind := bind_arg_l["value"].([]any)
				ident := ((bind[0].(map[string]any))["value"]).(string)
				bind_t, err := parse(bind[1].(map[string]any), env)
				if err != nil {

				}
				binds[i] = struct {
					string
					ast_t
					bool
				}{ident, bind_t, false}
			}

			for i := range binds {
				b := binds[i]
				old_t, bound := env[b.string]
				env[b.string] = b.ast_t
				if bound {
					binds[i].bool = true
					binds[i].ast_t = old_t
				}
			}

			body_t, err := parse(items[i+2].(map[string]any), env)
			if err != nil {
				return nil, err
			}

			// remap old env back
			for _, b := range binds {
				if b.bool {
					env[b.string] = b.ast_t
				} else {
					delete(env, b.string)
				}
			}

			return body_t, nil

		case "begin":
			expr_l := items[i+1:]
			var type_res ast_t
			var err error
			for _, expr := range expr_l {
				expr := expr.(map[string]any)
				type_res, err = parse(expr, env)
				if err != nil {
					return nil, err
				}

			}

			return type_res, nil

		case "set!":
			ident := items[i+1].(map[string]any)
			if ident["type"].(string) != "symbol" {
				return nil, fmt.Errorf("Non-symbol given as ident in set expr")
			}
			ident_str := ident["value"].(string)
			old_type, prs := env[ident_str]
			if !prs {
				return nil, fmt.Errorf("Cannot set undefined var `%s`", ident_str)
			}
			bound_val, err := parse(items[i+2].(map[string]any), env)
			if err != nil {
				return nil, err
			}
			if bound_val.Type() != old_type.Type() {
				return nil, fmt.Errorf("Cannot set %s, type %s to %s", ident_str, old_type, bound_val)
			}
			env[ident_str] = bound_val

			return ast_void{}, nil

		case "+":
			fallthrough
		case "-":
			fallthrough
		case "*":
			fallthrough
		case "/":
			return parseBinOp(items[i+1:], item_v, number_t, number_t, env)

		case "=":
			fallthrough
		case "<":
			return parseBinOp(items[i+1:], item_v, number_t, boolean_t, env)

		case "and":
			fallthrough
		case "or":
			return parseBinOp(items[i+1:], item_v, boolean_t, boolean_t, env)

		case "not":
			arg, err := parse(items[i+1].(map[string]any), env)
			if err != nil {
				return nil, err
			}
			if arg.Type() != boolean_t {
				return nil, fmt.Errorf("`(%s [boolean])` used with type [%s]", item_v, arg.Type())
			}

			return newLiteral(boolean_t)
		default:
			return nil, fmt.Errorf("unhandled symbol type '%s'", item_v)
		}
	}

	return ast, nil
}

func newLiteral(tipe string) (ast_t, error) {
	switch tipe {
	case "number":
		fallthrough
	case number_t:
		return ast_num{}, nil
	case "boolean":
		fallthrough
	case boolean_t:
		return ast_bool{}, nil
	default:
		return nil, fmt.Errorf("Invalid type used as literal: %s", tipe)
	}
}

func parseObjLit(sexp []any, env environ) (ast_t, error) {
	fields := map[string]ast_t{}
	for _, next := range sexp {
		next := next.(map[string]any)
		next_t := next["type"].(string)
		if next_t != "list" {
			break
		}

		field := next["value"].([]any)
		if len(field) != 2 {
			return nil, fmt.Errorf("Expected object field of form [ident expr]")
		}

		ident := field[0].(map[string]any)
		ident_t := ident["type"].(string)
		if ident_t != "symbol" {
			return nil, fmt.Errorf("Expected symbol in object field ident, got: %s", ident_t)
		}
		ident_str := ident["value"].(string)

		bound := field[1].(map[string]any)
		bound_v, err := parse(bound, env)
		if err != nil {
			return nil, err
		}

		if _, prs := fields[ident_str]; prs {
			return nil, fmt.Errorf("Object cannot have duplicate field `%s`", ident_str)
		}

		fields[ident_str] = bound_v

	}
	return ast_obj{fields}, nil
}

func parseBinOp(sexp []any, op, arg_t, ret_t string, env environ) (ast_t, error) {
	left, err := parse(sexp[0].(map[string]any), env)
	if err != nil {
		return nil, err
	}
	if left.Type() != arg_t {
		return nil, fmt.Errorf("`(%s [%s] %s)` used with type [%s]", op, arg_t, arg_t, left.Type())
	}

	right, err := parse(sexp[1].(map[string]any), env)
	if err != nil {
		return nil, err
	}
	if right.Type() != arg_t {
		return nil, fmt.Errorf("`(%s %s [%s])` used with type [%s]", op, arg_t, arg_t, right.Type())
	}

	return newLiteral(ret_t)
}
