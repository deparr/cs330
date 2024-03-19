package types

import "fmt"

func Parse(json map[string]any) (ast_t, error) {
	item_t := json["type"].(string)
	if item_t != "list" && item_t != "symbol" {
		return newLiteral(item_t)
	} else if item_t == "symbol" {
		// env lookup
		panic("TODO: env lookup for symbol sexp passed to types.Parse()")
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
		fmt.Printf("%d %v\n", i, item)
		// if type is not list ??
		// what to do here?
		// can't switch on value because that can be any type
		// WARN: this probably shouldn't be here
		item_v = item["value"].(string)
		switch item_v {
		case "object":
			ast, err = parseObjLit(items[i+1:])
			if err != nil {
				return nil, err
			}
			i++
			obj_len := len((ast.(ast_obj)).Fields)
			fmt.Println("NEXT AFTER OBJ", items[i+obj_len:])
			i += obj_len
		case "field":
			obj, err := Parse(items[i+1].(map[string]any))
			if err != nil {
				return nil, err
			}
			if obj.Type() != object_t {
				return nil, fmt.Errorf("'%s' type is not field-accessible", obj.Type())
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
			_ = (arg_l[0].(map[string].any))["value"].(string)
			arg_type, err := newLiteral(arg_l[2].(string))
			if err != nil {
				return nil, err
			}

			// env extend ??

			body_type, err := Parse(items[i+2].(map[string]any))
			if err != nil {
				return nil, err
			}
			return ast_func{arg_type, body_type}, nil

		case "+":
			fallthrough
		case "-":
			fallthrough
		case "*":
			fallthrough
		case "/":
			return parseBinOp(items[i+1:], item_v, number_t, number_t)

		case "=":
			fallthrough
		case "<":
			return parseBinOp(items[i+1:], item_v, number_t, boolean_t)

		case "and":
			fallthrough
		case "or":
			return parseBinOp(items[i+1:], item_v, boolean_t, boolean_t)

		case "not":
			arg, err := Parse(items[i+1].(map[string]any))
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
	case number_t:
		return ast_num{}, nil
	case boolean_t:
		return ast_bool{}, nil
	default:
		return nil, fmt.Errorf("Invalid type used as literal: %s", tipe)
	}
}

func parseObjLit(sexp []any) (ast_t, error) {
	fields := map[string]ast_t{}
	for _, next := range sexp {
		next := next.(map[string]any)
		next_t := next["type"].(string)
		// WARN: this should always be able to go to the end of the obj list
		//			though what happens when an invalid field is passed in ??
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
		bound_v, err := Parse(bound)
		// bound_v, err := newLiteral(bound["type"].(string))
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

func parseBinOp(sexp []any, op, arg_t, ret_t string) (ast_t, error) {
	left, err := Parse(sexp[0].(map[string]any))
	if err != nil {
		return nil, err
	}
	if left.Type() != arg_t {
		return nil, fmt.Errorf("`(%s [%s] %s)` used with type [%s]", op, arg_t, arg_t, left.Type())
	}

	right, err := Parse(sexp[1].(map[string]any))
	if err != nil {
		return nil, err
	}
	if right.Type() != arg_t {
		return nil, fmt.Errorf("`(%s %s [%s])` used with type [%s]", op, arg_t, arg_t, right.Type())
	}

	return newLiteral(ret_t)
}
