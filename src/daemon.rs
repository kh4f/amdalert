use std::{
    sync::Arc,
    sync::atomic::{AtomicBool, Ordering},
    thread,
    time::Duration,
};

use windows::{
    Win32::{
        Foundation::{CloseHandle, HANDLE},
        Storage::FileSystem::{PIPE_ACCESS_DUPLEX, ReadFile},
        System::Pipes::{
            ConnectNamedPipe, CreateNamedPipeW, DisconnectNamedPipe, PIPE_READMODE_MESSAGE,
            PIPE_TYPE_MESSAGE, PIPE_WAIT,
        },
    },
    core::{PCWSTR, w},
};

use crate::AnyResult;

pub const DAEMON_FLAG: &str = "--daemon";
pub const STOP_COMMAND: &str = "STOP";
const PIPE_NAME: PCWSTR = w!(r"\\.\pipe\amdlert");

pub fn run() -> AnyResult {
    let running = Arc::new(AtomicBool::new(true));
    let flag = Arc::clone(&running);

    let monitor = thread::spawn(move || {
        loop {
            if !flag.load(Ordering::Relaxed) {
                break;
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
