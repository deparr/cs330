# cs 330 typechecker
### David Parrott - dmparr22
---
# Building
This is a Go project, `Go` is required to build. You can get `go` here [go.dev](https://go.dev/dl/).

You should be fine with either `go1.21` or `go1.22`, I used `1.22`

# Typechecker
This is some of the worst code I've ever written.
The main file is `main.go`, run it with the following commands:
```sh
# run the type, expects json on stdin
go run main.go

racket stojson.rkt < s-exp.txt | go run main.go

# or if you want to run the binary directly
go build .
./tc_330
racket stojson.rkt < s-exp.txt | ./tc_330
./tc_330 < s-exp.json

# to remove build artifacts
go clean
```
