fn main() {
    cc::Build::new()
        .cpp(true)
        .files([
            "src/adlx_bridge.cpp",
            "ADLX/SDK/Platform/Windows/WinAPIs.cpp",
            "ADLX/SDK/ADLXHelper/Windows/Cpp/ADLXHelper.cpp",
        ])
        .compile("adlx_bridge");

    println!("cargo:rerun-if-changed=src/adlx_bridge.cpp");
}
