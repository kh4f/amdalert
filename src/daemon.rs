use std::{
    env,
    process::Command,
    sync::Arc,
    sync::atomic::{AtomicBool, Ordering},
    thread,
    time::Duration,
};

use windows::{
    Win32::{
        Foundation::{CloseHandle, GENERIC_READ, GENERIC_WRITE, HANDLE},
        Storage::FileSystem::{
            CreateFileW, FILE_ATTRIBUTE_NORMAL, FILE_SHARE_NONE, OPEN_EXISTING, PIPE_ACCESS_DUPLEX,
            ReadFile, WriteFile,
        },
        System::Console::FreeConsole,
        System::Pipes::{
            ConnectNamedPipe, CreateNamedPipeW, DisconnectNamedPipe, PIPE_READMODE_MESSAGE,
            PIPE_TYPE_MESSAGE, PIPE_WAIT, WaitNamedPipeW,
        },
        UI::WindowsAndMessaging::{MB_ICONWARNING, MB_OK, MessageBoxW},
    },
    core::{HSTRING, PCWSTR, w},
};

use crate::{AnyResult, adlx, config};

pub const DAEMON_FLAG: &str = "--daemon";
pub const STOP_COMMAND: &str = "STOP";
const PIPE_NAME: PCWSTR = w!(r"\\.\pipe\amdlert");

pub fn run() -> AnyResult {
    unsafe { FreeConsole() }.ok();

    let running = Arc::new(AtomicBool::new(true));
    let flag = Arc::clone(&running);

    let monitor = thread::spawn(move || {
        loop {
            if !flag.load(Ordering::Relaxed) {
                break;
            }

            if let Ok(gpu) = adlx::gpu_info() {
                let cfg = config::Config::load();
                if gpu.temperature > cfg.threshold {
                    show_warning(gpu.temperature, cfg.threshold);
                }
            }

            thread::park_timeout(Duration::from_secs(10));
        }
    });

    let pipe = create_pipe();
    loop {
        unsafe { ConnectNamedPipe(pipe, None) }?;
        let message = read_message(pipe)?;
        unsafe { DisconnectNamedPipe(pipe) }?;

        if message == STOP_COMMAND {
            break;
        }
    }

    running.store(false, Ordering::Relaxed);
    monitor.thread().unpark();
    monitor.join().ok();
    unsafe { CloseHandle(pipe) }?;
    Ok(())
}

pub fn spawn() -> AnyResult {
    Command::new(env::current_exe()?).arg(DAEMON_FLAG).spawn()?;
    Ok(())
}

pub fn is_daemon_running() -> bool {
    unsafe { WaitNamedPipeW(PIPE_NAME, 0) }.into()
}

pub fn send_message(msg: &str) -> AnyResult {
    let pipe = unsafe {
        CreateFileW(
            PIPE_NAME,
            (GENERIC_READ | GENERIC_WRITE).0,
            FILE_SHARE_NONE,
            None,
            OPEN_EXISTING,
            FILE_ATTRIBUTE_NORMAL,
            None,
        )
    }?;

    unsafe { WriteFile(pipe, Some(msg.as_bytes()), None, None) }?;
    unsafe { CloseHandle(pipe) }?;

    Ok(())
}

fn show_warning(temp: u32, threshold: u32) {
    let title = HSTRING::from("🚨 AMDlert");
    let text = HSTRING::from(format!("{temp}°C exceeds {threshold}°C threshold"));

    unsafe {
        MessageBoxW(None, &text, &title, MB_OK | MB_ICONWARNING);
    }
}

fn create_pipe() -> HANDLE {
    unsafe {
        CreateNamedPipeW(
            PIPE_NAME,
            PIPE_ACCESS_DUPLEX,
            PIPE_TYPE_MESSAGE | PIPE_READMODE_MESSAGE | PIPE_WAIT,
            1,
            64,
            64,
            0,
            None,
        )
    }
}

fn read_message(pipe: HANDLE) -> windows::core::Result<String> {
    let mut buffer = [0u8; 64];
    let mut bytes_read = 0;

    unsafe { ReadFile(pipe, Some(&mut buffer), Some(&mut bytes_read), None) }?;

    Ok(String::from_utf8_lossy(&buffer[..bytes_read as usize]).into_owned())
}
