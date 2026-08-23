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
        System::{
            Console::{
                ATTACH_PARENT_PROCESS, AllocConsole, AttachConsole, CONSOLE_MODE, COORD,
                ENABLE_VIRTUAL_TERMINAL_PROCESSING, GetConsoleMode, GetStdHandle, SMALL_RECT,
                STD_OUTPUT_HANDLE, SetConsoleMode, SetConsoleScreenBufferSize,
                SetConsoleWindowInfo,
            },
            Registry::{
                HKEY, HKEY_CURRENT_USER, KEY_QUERY_VALUE, KEY_SET_VALUE, REG_SZ, RegCloseKey,
                RegDeleteValueW, RegOpenKeyExW, RegQueryValueExW, RegSetValueExW,
            },
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

pub fn run() -> AnyResult {
    attach_console()?;

    let mut cfg = config::Config::load();
    let version = get_app_version()?;

    loop {
        let daemon_running = daemon::is_daemon_running();
        let gpu = adlx::gpu_info()?;

        print_menu(daemon_running, &gpu, cfg.threshold, &version)?;

        match read_choice()?.as_str() {
            "1" if daemon_running => {
                daemon::send_message(STOP_COMMAND)?;
                feedback("Daemon stopped")?;
            }
            "1" => {
                daemon::spawn()?;
                feedback("Daemon started")?;
            }
            "2" => {
                print!("New threshold (°C): ");
                io::stdout().flush()?;
                match read_choice()?.parse::<u32>() {
                    Ok(val) => {
                        cfg.set_threshold(val)?;
                        feedback(&format!("Threshold set to {val}°C"))?;
                    }
                    Err(_) => {
                        feedback("Invalid number")?;
                    }
                }
            }
            "3" if autostart_enabled() => {
                set_autostart(false)?;
                feedback("Autostart disabled")?;
            }
            "3" => {
                set_autostart(true)?;
                feedback("Autostart enabled")?;
            }
            "4" => break Ok(()),
            _ => {}
        }
    }
}

fn attach_console() -> AnyResult {
    if unsafe { AttachConsole(ATTACH_PARENT_PROCESS) }.is_err() {
        unsafe { AllocConsole() }?;

        let (cols, rows): (i16, i16) = (39, 14);
        let handle = unsafe { GetStdHandle(STD_OUTPUT_HANDLE) }?;

        #[expect(unused_must_use)]
        unsafe {
            SetConsoleWindowInfo(
                handle,
                true,
                &SMALL_RECT {
                    Left: 0,
                    Top: 0,
                    Right: cols - 1,
                    Bottom: rows - 1,
                },
            )
        };
        #[expect(unused_must_use)]
        unsafe {
            SetConsoleScreenBufferSize(handle, COORD { X: cols, Y: rows })
        };
    }

    // Enable ANSI escape codes
    let mut mode = CONSOLE_MODE(0);
    let handle = unsafe { GetStdHandle(STD_OUTPUT_HANDLE) }?;
    if unsafe { GetConsoleMode(handle, &mut mode) }.is_ok() {
        unsafe { SetConsoleMode(handle, mode | ENABLE_VIRTUAL_TERMINAL_PROCESSING) }?;
    }
    Ok(())
}

fn print_menu(
    daemon_running: bool,
    gpu: &adlx::GpuInfo,
    threshold: u32,
    version: &str,
) -> AnyResult {
    let (green, red, reset) = ("\x1B[32m", "\x1B[31m", "\x1B[0m");

    print!(
        "\x1B[2J\x1B[H\
AMDlert v{version}

{gpu_name}
  Temperature:  {gpu_temperature}°C
  Fan speed:    {gpu_fan_speed} RPM

Actions
  1) {daemon_action:<19} {daemon_status}
  2) {threshold_action:<19} {threshold}°C
  3) {autostart_action:<19} {autostart_status}
  4) Exit

> ",
        gpu_name = gpu.name,
        gpu_temperature = gpu.temperature,
        gpu_fan_speed = gpu.fan_speed,
        daemon_action = if daemon_running {
            "Stop daemon"
        } else {
            "Start daemon"
        },
        daemon_status = if daemon_running {
            format!("{green}running{reset}")
        } else {
            format!("{red}not running{reset}")
        },
        threshold_action = "Set threshold",
        autostart_action = if autostart_enabled() {
            "Disable autostart"
        } else {
            "Enable autostart"
        },
        autostart_status = if autostart_enabled() {
            format!("{green}enabled{reset}")
        } else {
            format!("{red}disabled{reset}")
        }
    );

    io::stdout().flush()?;
    Ok(())
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

fn feedback(msg: &str) -> AnyResult {
    print!("{msg}");
    io::stdout().flush()?;
    thread::sleep(Duration::from_millis(500));
    Ok(())
}
