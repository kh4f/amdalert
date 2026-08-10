# Changelog


## &ensp; [` 📦 v4.0.2  `](https://github.com/kh4f/amdlert/compare/v4.0.1...v4.0.2)

### &emsp; 🩹 Fixes
- **No console window on autostart**: the daemon now detaches from the console via `FreeConsole`, so starting it with autostart (or from the CLI) no longer flashes or keeps an empty terminal window. [🡥](https://github.com/kh4f/amdlert/commit/f8078cd)

##### &emsp;&emsp; [Commit log](https://github.com/kh4f/amdlert/compare/v4.0.1...v4.0.2) &ensp;•&ensp; Aug 10, 2026


## &ensp; [` 📦 v4.0.1  `](https://github.com/kh4f/amdlert/compare/v4.0.0...v4.0.1)

### &emsp; 🩹 Fixes
- **Daemon survives closing the terminal**: the daemon is now spawned detached from the console, so it keeps monitoring in the background even after the terminal window is closed, instead of being terminated with it. [🡥](https://github.com/kh4f/amdlert/commit/e4f298b)

### &emsp; 📦 Distribution
- **Release zip bundles an `AMDlert/` folder**: the archive now contains a ready-to-move folder with the executable inside, so the whole app can be dropped into `C:/Program Files/` as-is. [🡥](https://github.com/kh4f/amdlert/commit/7931bda)

### &emsp; 🎨 Style
- **Siren emoji in the warning title**: the overheat alert now opens with a 🚨 prefix for extra visibility. [🡥](https://github.com/kh4f/amdlert/commit/3b58795)

##### &emsp;&emsp; [Commit log](https://github.com/kh4f/amdlert/compare/v4.0.0...v4.0.1) &ensp;•&ensp; Aug 9, 2026


## &ensp; [` 📦 v4.0.0  `](https://github.com/kh4f/amdlert/compare/v3.0.2...v4.0.0)

### &emsp; 📢 BREAKING CHANGES
- **Complete rewrite in Rust**: the app has been fully rewritten from Go to Rust. It now ships as a single `AMDlert.exe` that covers both the interactive CLI and the background daemon — no separate daemon binary anymore. [🡥](https://github.com/kh4f/amdlert/commit/97ce17c)
- **Config format changed**: settings are now stored in `config.json` instead of the old `config.yml`. [🡥](https://github.com/kh4f/amdlert/commit/30c40c4)
- **Fan speed alerts removed**: monitoring now focuses on temperature only, keeping the app simpler. [🡥](https://github.com/kh4f/amdlert/commit/6256231)

### &emsp; ⚡ Performance
- **~32× smaller binary and faster startup**: the Rust rewrite combined with a size-optimized release profile (`opt-level = "z"`, LTO, single codegen unit) dramatically shrinks the executable and speeds it up. [🡥](https://github.com/kh4f/amdlert/commit/35419cf)

### &emsp; 🎨 Style
- **Cleaner CLI menu**: the interactive menu was made more concise. [🡥](https://github.com/kh4f/amdlert/commit/f506702)

##### &emsp;&emsp; [Commit log](https://github.com/kh4f/amdlert/compare/v3.0.2...v4.0.0) &ensp;•&ensp; Aug 9, 2026


## &ensp; [` 📦 v3.0.2  `](https://github.com/kh4f/amdalert/compare/v3.0.1...v3.0.2)

### &emsp; 🩹 Fixes
- **Reliable startup config loading**: the daemon now always reads `config.yml` from the executable directory, so custom temperature thresholds no longer get lost when `AMDAlertDaemon.exe` is launched from Windows startup with a different working directory. [🡥](https://github.com/kh4f/amdalert/commit/8c2c637)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/amdalert/compare/v3.0.1...v3.0.2) &ensp;•&ensp; Apr 6, 2026


## &ensp; [` 📦 v3.0.1  `](https://github.com/kh4f/amdalert/compare/v3.0.0...v3.0.1)

### &emsp; 🎨 Style
- **Sharper app icon**: fixed icon blurriness when displayed at small sizes. [🡥](https://github.com/kh4f/amdalert/commit/2a1558c)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/amdalert/compare/v3.0.0...v3.0.1) &ensp;•&ensp; Mar 25, 2026


## &ensp; [` 📦 v3.0.0  `](https://github.com/kh4f/amdalert/compare/v2.0.1...v3.0.0)

### &emsp; 🧨 BREAKING CHANGES
- **YAML Config Migration**: configuration is now stored in `config.yml` using YAML syntax instead of `settings.json`. [🡥](https://github.com/kh4f/amdalert/commit/a856e4e)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/amdalert/compare/v2.0.1...v3.0.0) &ensp;•&ensp; Mar 25, 2026


## &ensp; [` 📦 v2.0.1  `](https://github.com/kh4f/amdalert/compare/v2.0.0...v2.0.1)

### &emsp; 🎨 Style
- **Updated app icon**: the application icon has been refreshed with a gradient. [🡥](https://github.com/kh4f/amdalert/commit/97e10b8)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/amdalert/compare/v2.0.0...v2.0.1) &ensp;•&ensp; Mar 23, 2026


## &ensp; [` 📦 v2.0.0  `](https://github.com/kh4f/amdalert/compare/v1.0.0...v2.0.0)

### &emsp; 🧨 BREAKING CHANGES
- **Renamed config fields**: configuration fields are now `alertTemp` and `fanOffAlertTemp`. [🡥](https://github.com/kh4f/amdalert/commit/7906d1f)
- **Config file renamed**: configuration is now loaded from `settings.json` instead of `config.json`. [🡥](https://github.com/kh4f/amdalert/commit/fef3986)
- **Split into CLI and daemon**: the application now consists of two executables: `AMDAlert.exe` for CLI and `AMDAlertDaemon.exe` for the background service. [🡥](https://github.com/kh4f/amdalert/commit/e3a6b05)

### &emsp; 📚 Documentation
- **Enriched documentation**: added project overview, feature list, installation and usage instructions, logo, and demo screenshot. [🡥](https://github.com/kh4f/amdalert/commit/9b965d5)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/amdalert/compare/v1.0.0...v2.0.0) &ensp;•&ensp; Mar 22, 2026


## &ensp; [` 📦 v1.0.0  `](https://github.com/kh4f/amdalert/compare/v0.3.0...v1.0.0)

### &emsp; 🧨 BREAKING CHANGES
- **Project renamed to `AMDAlert`**: all code, modules, binaries, registry keys, and documentation now use the `AMDAlert` name instead of `amdalert`. [🡥](https://github.com/kh4f/amdalert/commit/05c825d)

### &emsp; 🎁 Features
- **Windows autostart support**: the app can now be added to or removed from Windows startup via the CLI. [🡥](https://github.com/kh4f/amdalert/commit/3bce049)
- **Application icon**: a custom `icon.ico` is now embedded in the Windows executable. [🡥](https://github.com/kh4f/amdalert/commit/3413dac)

### &emsp; 🎨 Style
- **Improved CLI layout**: the CLI output is now more readable, with grouped status and actions, and clearer alert information. [🡥](https://github.com/kh4f/amdalert/commit/8f00e7f)
- **Version info in CLI banner**: the CLI header now displays the current version for better traceability. [🡥](https://github.com/kh4f/amdalert/commit/407175f)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/amdalert/compare/v0.3.0...v1.0.0) &ensp;•&ensp; Mar 16, 2026


## &ensp; [` 📦 v0.3.0  `](https://github.com/kh4f/amdalert/compare/v0.2.0...v0.3.0)

### &emsp; 🎁 Features
- **Interactive CLI menu**: added daemon control and GPU information display directly in the command-line interface. [🡥](https://github.com/kh4f/amdalert/commit/4daea2a)
- **Custom temperature limits**: configuration now supports `maxTemp` and `maxFanOffTemp` settings loaded from `config.json`. [🡥](https://github.com/kh4f/amdalert/commit/b781a3f)
- **Configuration management**: daemon now monitors external changes to `config.json` and reloads automatically during runtime. [🡥](https://github.com/kh4f/amdalert/commit/0f4aa28)
- **Temperature threshold configuration**: added interactive prompts to set maximum temperature and maximum fan-off temperature from the CLI. [🡥](https://github.com/kh4f/amdalert/commit/a3fdfb3)
- **Configuration visibility**: current temperature and fan-off temperature settings now display in the CLI interface. [🡥](https://github.com/kh4f/amdalert/commit/5583305)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/amdalert/compare/v0.2.0...v0.3.0) &ensp;•&ensp; Mar 14, 2026


## &ensp; [` 📦 v0.2.0  `](https://github.com/kh4f/amdalert/compare/v0.1.0...v0.2.0)

### &emsp; 🎁 Features
- **GPU fan and temperature alerts**: the application now shows alerts if GPU temperature exceeds 40°C and the fan is not spinning, or if temperature exceeds 60°C. [🡥](https://github.com/kh4f/amdalert/commit/f7e0377)
- **Display GPU fan speed**: added support for reading and printing GPU fan speed (RPM) every 10 seconds together with temperature. [🡥](https://github.com/kh4f/amdalert/commit/919290d)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/amdalert/compare/v0.1.0...v0.2.0) &ensp;•&ensp; Mar 11, 2026


## &ensp; [` 📦 v0.1.0  `](https://github.com/kh4f/amdalert/commits/v0.1.0)

### &emsp; 🎁 Features
- **Periodic GPU temperature display**: the application now prints GPU temperature every 10 seconds for continuous monitoring. [🡥](https://github.com/kh4f/amdalert/commit/342fcc9)
- **GPU temperature reading via ADL**: added AMD Display Library integration to read GPU temperature. [🡥](https://github.com/kh4f/amdalert/commit/3d97292)

### &emsp; ⚙️ Internal
- **Migrated to Go implementation**: replaced the previous Rust codebase with a Go implementation, including new project structure and ADL integration. [🡥](https://github.com/kh4f/amdalert/commit/ee53829)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/amdalert/commits/v0.1.0) &ensp;•&ensp; Mar 10, 2026