fn main() -> Result<(), Box<dyn std::error::Error>> {
    cc::Build::new()
        .cpp(true)
        .files([
            "src/adlx_bridge.cpp",
            "ADLX/SDK/Platform/Windows/WinAPIs.cpp",
            "ADLX/SDK/ADLXHelper/Windows/Cpp/ADLXHelper.cpp",
        ])
        .try_compile("adlx_bridge")?;

    winresource::WindowsResource::new()
        .set("FileDescription", "AMDlert")
        .set("ProductName", "AMDlert")
        .set("LegalCopyright", "Copyright © 2026 kh4f")
        .set_language(0x0409)
        .compile()?;

    println!("cargo:rerun-if-changed=src/adlx_bridge.cpp");
    Ok(())
}
