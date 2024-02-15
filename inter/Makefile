run: build
	@echo "running bind, expecting ast on stdin..."
	@./target/debug/bind

run-acorn: build
	@echo "running bind..."
	@acorn --ecma2024 | ./target/debug/bind


build:
	cargo build
	@printf "\ndone: binaries available in ./target/debug/\n\n"

test:
	cargo test

clean:
	cargo clean
