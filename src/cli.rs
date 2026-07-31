use std::io::{self, Write};

use crate::AnyResult;

pub fn run() -> AnyResult {
    loop {
        print_menu();

        match read_choice()?.as_str() {
            "1" => break Ok(()),
            _ => {}
        }
    }
}

fn print_menu() {
    print!(
        "\x1B[2J\x1B[H\
🚨 AMDlert

Actions
  1) Exit

> "
    );

    io::stdout().flush().ok();
}

fn read_choice() -> AnyResult<String> {
    let mut input = String::new();
    io::stdin().read_line(&mut input)?;
    Ok(input.trim().to_string())
}
