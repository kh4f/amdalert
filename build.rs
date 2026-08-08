use std::{
    env,
    error::Error,
    path::{Path, PathBuf},
};

fn main() -> Result<(), Box<dyn Error>> {
    cc::Build::new()
        .cpp(true)
        .files([
            "src/adlx_bridge.cpp",
            "ADLX/SDK/Platform/Windows/WinAPIs.cpp",
            "ADLX/SDK/ADLXHelper/Windows/Cpp/ADLXHelper.cpp",
        ])
        .try_compile("adlx_bridge")?;

    let ico_path = PathBuf::from(env::var("OUT_DIR")?).join("icon.ico");
    let ico_path = ico_path
        .to_str()
        .ok_or("failed to convert ico path to &str")?;
    svico::convert(Path::new("assets/icon.svg"), ico_path, &[16, 24, 32, 256])?;

    winresource::WindowsResource::new()
        .set("FileDescription", "AMDlert")
        .set("ProductName", "AMDlert")
        .set("LegalCopyright", "Copyright © 2026 kh4f")
        .set_language(0x0409)
        .set_icon(ico_path)
        .compile()?;

    println!("cargo:rerun-if-changed=src/adlx_bridge.cpp");
    println!("cargo:rerun-if-changed=assets/icon.svg");
    Ok(())
}
