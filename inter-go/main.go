package main

import (
	"encoding/json"
	"fmt"
	"dmparr22/inter_330/ast"
	"io"
	"os"
)

func main() {
	jsonBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read from stdin: %s\n", err)
		os.Exit(1)
	}

	var astMap = make(map[string]interface{})
	err = json.Unmarshal(jsonBytes, &astMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to unmarshal json: %s\n", err)
		os.Exit(1)
	}

	ast, err := ast.New(astMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create ast: %s\n", err)
		os.Exit(3)
	}

	fmt.Println(ast)
	result, err := ast.Eval()
	if err != nil {
		fmt.Printf("(error \"%s banana\")\n", err)
		return
	}

	fmt.Printf("(value %s)\n", result)
}
