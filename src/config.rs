use std::{env, fs, path::PathBuf};

use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct Config {
    pub threshold: u32,
}

impl Default for Config {
    fn default() -> Self {
        Self { threshold: 50 }
    }
}

impl Config {
    fn path() -> PathBuf {
        let mut path = env::current_exe().expect("failed to get exe path");
        path.pop();
        path.push("config.json");
        path
    }

    pub fn load() -> Self {
        let path = Self::path();
        fs::read_to_string(&path)
            .ok()
            .and_then(|s| serde_json::from_str(&s).ok())
            .unwrap_or_default()
    }

    pub fn save(&self) -> Result<(), String> {
        let json = serde_json::to_string_pretty(self).map_err(|e| e.to_string())?;
        fs::write(Self::path(), json).map_err(|e| e.to_string())
    }

    pub fn set_threshold(&mut self, value: u32) -> Result<(), String> {
        self.threshold = value;
        self.save()
    }
}
