use std::{
    env,
    ffi::c_void,
    io::{self, Write},
    iter, ptr, slice, thread,
    time::Duration,
};

use windows::{
    Win32::{
        Storage::FileSystem::{
            GetFileVersionInfoSizeW, GetFileVersionInfoW, VS_FIXEDFILEINFO, VerQueryValueW,
        },
        System::Registry::{
            HKEY, HKEY_CURRENT_USER, KEY_QUERY_VALUE, KEY_SET_VALUE, REG_SZ, RegCloseKey,
            RegDeleteValueW, RegOpenKeyExW, RegQueryValueExW, RegSetValueExW,
        },
    },
    core::{PCWSTR, w},
};

use crate::{
    AnyResult, adlx, config,
    daemon::{self, STOP_COMMAND},
};

const APP_NAME: PCWSTR = w!("AMDlert");
const RUN_KEY: PCWSTR = w!(r"Software\Microsoft\Windows\CurrentVersion\Run");
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
            "3" if autostart_enabled() => {
                set_autostart(false)?;
                println!("Autostart disabled");
                thread::sleep(FEEDBACK_DELAY);
            }
            "3" => {
                set_autostart(true)?;
                println!("Autostart enabled");
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
  1) {daemon_action:<19} {daemon_status}
  2) {threshold_action:<22} {threshold}°C
  3) {autostart_action:<19} {autostart_status}
  4) Exit

> ",
        gpu_name = gpu.name,
        gpu_temperature = gpu.temperature,
        gpu_fan_speed = gpu.fan_speed,
        daemon_status = if daemon_running { "🟢 running" } else { "🔴 not running" },
        autostart_status = if autostart_enabled() { "🟢 enabled" } else { "🔴 disabled" },
        daemon_action = if daemon_running { "Stop daemon" } else { "Run daemon" },
        threshold_action = "Set threshold",
        autostart_action = if autostart_enabled() { "Disable autostart" } else { "Enable autostart" },
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

fn autostart_enabled() -> bool {
    let mut hkey = HKEY::default();
    if unsafe {
        RegOpenKeyExW(
            HKEY_CURRENT_USER,
            RUN_KEY,
            Some(0),
            KEY_QUERY_VALUE,
            &mut hkey,
        )
    }
    .is_err()
    {
        return false;
    }

    let exists = unsafe { RegQueryValueExW(hkey, APP_NAME, None, None, None, None).is_ok() };

    #[expect(unused_must_use)]
    unsafe {
        RegCloseKey(hkey)
    };
    exists
}

fn set_autostart(enabled: bool) -> AnyResult {
    let mut hkey = HKEY::default();
    if unsafe {
        RegOpenKeyExW(
            HKEY_CURRENT_USER,
            RUN_KEY,
            Some(0),
            KEY_SET_VALUE,
            &mut hkey,
        )
    }
    .is_err()
    {
        return Err("failed to open registry key".into());
    }

    if enabled {
        let exe = env::current_exe()?;
        let path = format!("\"{}\" --daemon", exe.display());
        let wide: Vec<u16> = path.encode_utf16().chain(iter::once(0)).collect();
        let data = unsafe { slice::from_raw_parts(wide.as_ptr() as *const u8, wide.len() * 2) };

        if unsafe { RegSetValueExW(hkey, APP_NAME, None, REG_SZ, Some(data)) }.is_err() {
            #[expect(unused_must_use)]
            unsafe {
                RegCloseKey(hkey)
            };
            return Err("failed to set registry value".into());
        }
    } else {
        #[expect(unused_must_use)]
        unsafe {
            RegDeleteValueW(hkey, APP_NAME)
        };
    }

    #[expect(unused_must_use)]
    unsafe {
        RegCloseKey(hkey)
    };
    Ok(())
}
