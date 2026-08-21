lint:
	cargo fmt --check
	cargo clippy -- -D warnings

build action="":
	{{ if action == "kill" { "taskkill //f //im AMDlert.exe || true" } else { "" } }}
	cargo build --release

refresh-icons:
	taskkill //f //im explorer.exe
	start explorer.exe