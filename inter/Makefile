run: build
	@echo "running eval, expecting ast on stdin..."
	cargo run --bin eval

run-acorn: build
	@echo "running eval..."
	acorn --ecma2024 | cargo run --bin eval


build:
	cargo build
	@printf "\ndone: binaries available in ./target/debug/\n\n"

test:
	cargo test

clean:
	cargo clean
