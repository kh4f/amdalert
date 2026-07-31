use std::ffi::{CStr, c_int};

mod ffi {
    use std::ffi::{c_char, c_int};

    unsafe extern "C" {
        pub fn init() -> c_int;
        pub fn gpu_info(
            name_buf: *mut c_char,
            buf_size: c_int,
            out_temp: *mut c_int,
            out_fan_speed: *mut c_int,
        ) -> c_int;
    }
}

macro_rules! call {
    ($fn:expr, $msg:expr) => {
        let code = unsafe { $fn };
        if code != 0 {
            return Err(format!("{}: {code}", $msg));
        }
    };
}

pub struct GpuInfo {
    pub name: String,
    pub temperature: u32,
    pub fan_speed: u32,
}

pub fn init() -> Result<(), String> {
    call!(ffi::init(), "ADLX init failed");
    Ok(())
}

pub fn gpu_info() -> Result<GpuInfo, String> {
    let mut name_buf = [0i8; 64];
    let mut temp: c_int = 0;
    let mut fan_speed: c_int = 0;

    call!(
        ffi::gpu_info(
            name_buf.as_mut_ptr(),
            name_buf.len() as c_int,
            &mut temp,
            &mut fan_speed,
        ),
        "failed to get GPU info"
    );

    let name = unsafe { CStr::from_ptr(name_buf.as_ptr()) }
        .to_string_lossy()
        .into_owned();

    Ok(GpuInfo {
        name,
        temperature: temp as u32,
        fan_speed: fan_speed as u32,
    })
}
