mod adlx;
mod cli;
mod daemon;

use std::{env, error::Error};

use daemon::DAEMON_FLAG;

pub type AnyResult<T = ()> = Result<T, Box<dyn Error>>;

fn main() -> AnyResult {
    adlx::init()?;

    if env::args().any(|arg| arg == DAEMON_FLAG) {
        daemon::run()
    } else {
        cli::run()
    }
}
