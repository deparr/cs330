# cs 330 interpreter -- loops
### David Parrott - dmparr22
---
# Building
I used the starter plait implementation, so only `racket` is required to run the interpreter.

# Evaluator with loops
```sh
# run the evaluator, expects ast input on stdin
racket inter.rkt

# or pipe in the ast file
racket inter.rkt < ast.json
acorn --ecma2024 | racket inter.rkt
```

