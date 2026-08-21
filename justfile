lint:
	cargo fmt --check
	cargo clippy -- -D warnings

build:
	cargo build --release

refresh-icons:
	taskkill //f //im explorer.exe
	start explorer.exe