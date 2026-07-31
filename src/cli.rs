use std::{
    io::{self, Write},
    thread,
    time::Duration,
};

use crate::{AnyResult, adlx, config};

const FEEDBACK_DELAY: Duration = Duration::from_millis(500);

pub fn run() -> AnyResult {
    let mut cfg = config::Config::load();

    loop {
        let gpu = adlx::gpu_info()?;

        print_menu(&gpu, cfg.threshold);

        match read_choice()?.as_str() {
            "1" => break Ok(()),
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
            _ => {}
        }
    }
}

fn print_menu(gpu: &adlx::GpuInfo, threshold: u32) {
    print!(
        "\x1B[2J\x1B[H\
🚨 AMDlert

GPU ({gpu_name})
  Temperature:  {gpu_temperature}°C
  Fan speed:    {gpu_fan_speed} RPM

Actions
  1) Exit
  2) Set threshold {threshold}°C

> ",
        gpu_name = gpu.name,
        gpu_temperature = gpu.temperature,
        gpu_fan_speed = gpu.fan_speed,
    );

    io::stdout().flush().ok();
}

fn read_choice() -> AnyResult<String> {
    let mut input = String::new();
    io::stdin().read_line(&mut input)?;
    Ok(input.trim().to_string())
}
