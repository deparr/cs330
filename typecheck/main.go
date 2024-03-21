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

	_type, err := types.Parse(astMap)
	if err != nil {
		fmt.Printf("(error \"%s (banana)\")\n", err.Error())
		return
	}

	fmt.Printf("(type %s)\n",  _type)
}
