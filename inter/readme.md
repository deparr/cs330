# cs 330 interpreter
### David Parrott - dmparr22
---
# Building
This is a rust project, so building with `cargo`, the Rust package manager, is preferred. You can install both `cargo` and `rustc` at [rustup.rs](https://rustup.rs).  

You *should* be fine with any recent-ish rust version, but if you want to be sure, I used version `stable 1.75`.

# Evaluator
The expression data structure is defined in `src/ast.rs`.  
The parser main file is `src/bin/bind.rs`, run it with the following commands:
```sh
# run the evaluator, expects ast input on stdin
make run

# pipe in ast file
cat ast.json |  make run

# use `acorn` to generate ast, assumes `acorn` is executable
make run-acorn

# or if you want to run the binary directly
make build
./target/debug/bind
acorn --ecma2024 | ./target/debug/bind
cat ast.json | ./target/debug/bind

# to remove build artifacts
make clean
```
This will also download and compile the necessary dependencies: [serde](https://github.com/serde-rs/serde) and [serde_json](https://github.com/serde-rs/json).
