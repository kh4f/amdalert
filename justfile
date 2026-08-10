lint:
	cargo fmt --check
	cargo clippy -- -D warnings

build:
	cargo build --release