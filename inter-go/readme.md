# cs 330 interpreter
### David Parrott - dmparr22
---
# Building
This is a Go project, `Go` is required to build. You can get `go` here [go.dev](https://go.dev/dl/).

You *should* be fine with either `go1.21` or `go1.22`, I used `1.22`

# Evaluator with functions
The expression data structure is defined in `ast/ast.go` (note: it's a little messy, sorry)
The main file is `main.go`, run it with the following commands:
```sh
# run the evaluator, expects ast input on stdin
go run main.go

# pipe in ast file
cat ast.json |  go run main.go

# or if you want to run the binary directly
go build .
./inter_330
acorn --ecma2024 | ./inter_330
cat ast.json | ./inter_330

# to remove build artifacts
go clean
```
