use std::io::{self, Write};

use crate::{AnyResult, adlx};

pub fn run() -> AnyResult {
    loop {
        let gpu = adlx::gpu_info()?;

        print_menu(&gpu);

        match read_choice()?.as_str() {
            "1" => break Ok(()),
            _ => {}
        }
    }
}

fn print_menu(gpu: &adlx::GpuInfo) {
    print!(
        "\x1B[2J\x1B[H\
🚨 AMDlert

GPU ({gpu_name})
  Temperature:  {gpu_temperature}°C
  Fan speed:    {gpu_fan_speed} RPM

Actions
  1) Exit

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
