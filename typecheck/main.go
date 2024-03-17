package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"dmparr22/tc_330/types"
)

func main() {
	jsonBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read from stdin: %s\n", err)
		os.Exit(1)
	}

	var astMap = make(map[string]any)
	err = json.Unmarshal(jsonBytes, &astMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to unmarshal json: %s\n", err)
		os.Exit(1)
	}

	parser := types.NewParser()
	parsedAst, err := parser.Parse(astMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(parsedAst)

	_type, err := types.Check(astMap)
	if err != nil {
		fmt.Printf("(error \"%s\")\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(_type)
}
