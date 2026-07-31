use std::{
    io::{self, Write},
    thread,
    time::Duration,
};

use crate::{
    AnyResult, adlx, config,
    daemon::{self, STOP_COMMAND},
};

const FEEDBACK_DELAY: Duration = Duration::from_millis(500);

pub fn run() -> AnyResult {
    let mut cfg = config::Config::load();

    loop {
        let daemon_running = daemon::is_daemon_running();
        let gpu = adlx::gpu_info()?;

        print_menu(daemon_running, &gpu, cfg.threshold);

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

fn print_menu(daemon_running: bool, gpu: &adlx::GpuInfo, threshold: u32) {
    print!(
        "\x1B[2J\x1B[H\
🚨 AMDlert

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
