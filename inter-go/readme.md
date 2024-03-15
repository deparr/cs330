# cs 330 interpreter
### David Parrott - dmparr22
---
# Building
This is a poorly written Go project, `Go` is required to build. You can get `go` here [go.dev](https://go.dev/dl/).

You *should* be fine with either `go1.21` or `go1.22`, I used `1.22`

# Evaluator with functions
The main file is `main.go`, run it with the following commands:
```sh
# run the evaluator, expects ast input on stdin
go run main.go

# pipe in ast file
go run main.go < ast.json

# or if you want to run the binary directly
go build .
./inter_330
acorn --ecma2024 | ./inter_330
./inter_330 < ast.json

# to remove build artifacts
go clean
```
