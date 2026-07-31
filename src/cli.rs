use std::{
    env,
    ffi::c_void,
    io::{self, Write},
    iter, ptr, thread,
    time::Duration,
};

use windows::{
    Win32::Storage::FileSystem::{
        GetFileVersionInfoSizeW, GetFileVersionInfoW, VS_FIXEDFILEINFO, VerQueryValueW,
    },
    core::{PCWSTR, w},
};

use crate::{
    AnyResult, adlx, config,
    daemon::{self, STOP_COMMAND},
};

const FEEDBACK_DELAY: Duration = Duration::from_millis(500);

pub fn run() -> AnyResult {
    let mut cfg = config::Config::load();
    let version = get_app_version()?;

    loop {
        let daemon_running = daemon::is_daemon_running();
        let gpu = adlx::gpu_info()?;

        print_menu(daemon_running, &gpu, cfg.threshold, &version);

        match read_choice()?.as_str() {
            "1" if daemon_running => {
                daemon::send_message(STOP_COMMAND)?;
                println!("Daemon stopped");
                thread::sleep(FEEDBACK_DELAY);
            }
            "1" => {
                daemon::spawn()?;
                println!("Daemon started");
                thread::sleep(FEEDBACK_DELAY);
            }
            "2" => {
                print!("New threshold (°C): ");
                io::stdout().flush().ok();
                match read_choice()?.parse::<u32>() {
                    Ok(val) => {
                        cfg.set_threshold(val)?;
                        println!("Threshold set to {val}°C");
                    }
                    Err(_) => {
                        println!("Invalid number");
                    }
                }
                thread::sleep(FEEDBACK_DELAY);
            }
            "4" => break Ok(()),
            _ => {}
        }
    }
}

fn print_menu(daemon_running: bool, gpu: &adlx::GpuInfo, threshold: u32, version: &str) {
    print!(
        "\x1B[2J\x1B[H\
🚨 AMDlert v{version}
  Homepage: https://github.com/kh4f/amdlert

GPU ({gpu_name})
  Temperature:  {gpu_temperature}°C
  Fan speed:    {gpu_fan_speed} RPM

Actions
  1) {daemon_action} {daemon_status}
  2) Set threshold {threshold}°C
  4) Exit

> ",
        gpu_name = gpu.name,
        gpu_temperature = gpu.temperature,
        gpu_fan_speed = gpu.fan_speed,
        daemon_action = if daemon_running {
            "Stop daemon"
        } else {
            "Run daemon"
        },
        daemon_status = if daemon_running {
            "🟢 running"
        } else {
            "🔴 not running"
        },
    );

    io::stdout().flush().ok();
}

fn read_choice() -> AnyResult<String> {
    let mut input = String::new();
    io::stdin().read_line(&mut input)?;
    Ok(input.trim().to_string())
}

fn get_app_version() -> AnyResult<String> {
    use std::os::windows::ffi::OsStrExt;

    let exe = env::current_exe()?;
    let wide: Vec<u16> = exe.as_os_str().encode_wide().chain(iter::once(0)).collect();

    let size = unsafe { GetFileVersionInfoSizeW(PCWSTR(wide.as_ptr()), None) };
    if size == 0 {
        return Err("failed to get file version info size".into());
    }

    let mut buf = vec![0u8; size as usize];
    unsafe {
        GetFileVersionInfoW(
            PCWSTR(wide.as_ptr()),
            None,
            size,
            buf.as_mut_ptr() as *mut c_void,
        )
    }?;

    let mut info: *mut VS_FIXEDFILEINFO = ptr::null_mut();
    let mut len: u32 = 0;
    let sub_block = w!("\\");

    if !unsafe {
        VerQueryValueW(
            buf.as_ptr() as *const c_void,
            sub_block,
            &mut info as *mut *mut VS_FIXEDFILEINFO as *mut *mut c_void,
            &mut len,
        )
    }
    .as_bool()
    {
        return Err("failed to query version info".into());
    }

    let verinfo = unsafe { *info };
    let major = (verinfo.dwProductVersionMS >> 16) as u16;
    let minor = (verinfo.dwProductVersionMS & 0xFFFF) as u16;
    let patch = (verinfo.dwProductVersionLS >> 16) as u16;
    Ok(format!("{major}.{minor}.{patch}"))
}
