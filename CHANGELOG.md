# Changelog


## &ensp; [` 📦 v0.2.0  `](https://github.com/kh4f/hotamd/compare/v0.1.0...v0.2.0)

### &emsp; 🎁 Features
- **GPU fan and temperature alerts**: the application now shows alerts if GPU temperature exceeds 40°C and the fan is not spinning, or if temperature exceeds 60°C. [🡥](https://github.com/kh4f/hotamd/commit/f7e0377)
- **Display GPU fan speed**: added support for reading and printing GPU fan speed (RPM) every 10 seconds together with temperature. [🡥](https://github.com/kh4f/hotamd/commit/919290d)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/hotamd/compare/v0.1.0...v0.2.0) &ensp;•&ensp; Mar 11, 2026


## &ensp; [` 📦 v0.1.0  `](https://github.com/kh4f/hotamd/commits/v0.1.0)

### &emsp; 🎁 Features
- **Periodic GPU temperature display**: the application now prints GPU temperature every 10 seconds for continuous monitoring. [🡥](https://github.com/kh4f/hotamd/commit/342fcc9)
- **GPU temperature reading via ADL**: added AMD Display Library integration to read GPU temperature. [🡥](https://github.com/kh4f/hotamd/commit/3d97292)

### &emsp; ⚙️ Internal
- **Migrated to Go implementation**: replaced the previous Rust codebase with a Go implementation, including new project structure and ADL integration. [🡥](https://github.com/kh4f/hotamd/commit/ee53829)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/hotamd/commits/v0.1.0) &ensp;•&ensp; Mar 10, 2026