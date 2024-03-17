package types

import "fmt"

type Parser struct {
	cur int
}

func NewParser() Parser {
	return Parser{0}
}

func (p *Parser) Parse(json map[string]any) (ast_t, error) {
	item_t := json["type"].(string)
	if item_t != "list" {
		return newNumBoolLiteral(item_t)
	}

	items := json["value"].([]any)

	for p.cur < len(items) {
		item := items[p.cur].(map[string]any)
		item_t = item["type"].(string)
		fmt.Printf("%d %v\n", p.cur, item)
		if item_t != "symbol" && item_t != "list" {
			return newNumBoolLiteral(item_t)
		}
		p.cur++
	}

	return nil, nil
}

func newNumBoolLiteral(tipe string) (ast_t, error) {
	switch tipe {
	case "number":
		return ast_num{}, nil
	case "boolean":
		return ast_bool{}, nil
	default:
		return nil, fmt.Errorf("Inavlid type used as literal: %s", tipe)
	}
}
