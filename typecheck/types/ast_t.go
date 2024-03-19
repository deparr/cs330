package types

import (
	"fmt"
	"strings"
)

type ast_t interface {
	Ast_t_f()
	Type() string
}

const (
	number_t  = "num"
	boolean_t = "bool"
	void_t    = "void"
	object_t  = "object"   // ??
	func_t    = "function" // ??
)

type ast_num struct{}
type ast_bool struct{}
type ast_void struct{}
type ast_func struct {
	Arg ast_t
	Ret ast_t
}
type ast_obj struct {
	Fields map[string]ast_t
}

func (t ast_num) Ast_t_f()  {}
func (t ast_bool) Ast_t_f() {}
func (t ast_void) Ast_t_f() {}
func (t ast_func) Ast_t_f() {}
func (t ast_obj) Ast_t_f()  {}

func (t ast_num) Type() string  { return number_t }
func (t ast_bool) Type() string { return boolean_t }
func (t ast_void) Type() string { return void_t }
func (t ast_func) Type() string { return func_t }
func (t ast_obj) Type() string  { return object_t }

func (t ast_num) String() string {
	return "(number)"
}

func (t ast_bool) String() string {
	return "(boolean)"
}

func (t ast_func) String() string {
	return fmt.Sprintf("(-> %s %s)", t.Arg, t.Ret)
}

func (t ast_obj) String() string {
	fields := []string{}
	for k, v := range t.Fields {
		fields = append(fields, fmt.Sprintf("[%s %s]", k, v))
	}

	return fmt.Sprintf("(object %s)", strings.Join(fields, " "))
}
