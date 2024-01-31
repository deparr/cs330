# cs 330 interpreter
### David Parrott - dmparr22
---
# Building
This is a rust project, so building with `cargo`, the Rust package manager, is preferred. You can install both `cargo` and `rustc` at [rustup.rs](https://rustup.rs).  

You *should* be fine with any recent-ish rust version, but if you want to be sure, I used version `stable 1.75`.

# Parser
The expression data structure is defined in `src/ast.rs`.  
The parser main file is `src/bin/parser.rs`, run it with  
```sh
# by default it expects the ast on stdin
cargo run --bin parser

# pipe in ast file
cat ast.json | cargo run --bin parser

# have the parser exec acorn and capture output, acorn needs to be in PATH
cargo run --bin parser -- --exec

# or if you want to run the binary directly
cargo build
./target/debug/parser
```
This will also download and compile the necessary dependencies: [serde](https://github.com/serde-rs/serde) and [serde_json](https://github.com/serde-rs/json).
