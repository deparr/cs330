package types

import (
	"fmt"
	"strings"
)

type ast_t interface {
	Ast_t_f()
}

type ast_num struct{}
type ast_bool struct{}
type ast_void struct{}
type ast_func struct {
	arg ast_t
	ret ast_t
}
type ast_obj struct {
	fields map[string]ast_t
}


func (t ast_num) Ast_t_f()  {}
func (t ast_bool) Ast_t_f() {}
func (t ast_void) Ast_t_f() {}
func (t ast_func) Ast_t_f() {}
func (t ast_obj) Ast_t_f()  {}

func (t ast_num) String() string {
	return "(number)"
}

func (t ast_bool) String() string {
	return "(boolean)"
}

func (t ast_func) String() string {
	return fmt.Sprintf("(-> %s %s)")
}

func (t ast_obj) String() string {
	fields := []string{}
	for k, v := range t.fields {
		fields = append(fields, fmt.Sprintf("[%s %s]", k, v))
	}

	return fmt.Sprintf("(obj %s)", strings.Join(fields, " "))
}
