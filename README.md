<div align="center">
	<img alt="logo" src="assets/logo.png">
	<br>
	A lightweight Windows utility that <b>alerts you when your AMD GPU overheats</b><br>Perfect for older cards like the RX 550 where the fan may silently stop 🫠
	<br><br>
	<b>
		<a href="#-features">Features</a>&nbsp; •&nbsp;
		<a href="#-installation">Installation</a>&nbsp; •&nbsp;
		<a href="#%EF%B8%8F-usage">Usage</a>
	</b>
	<br><br>
	<img alt="demo" src="assets/demo.png">
</div>

## 🔥 Features

- Monitors GPU state and triggers alerts on overheating
- Detects fan failures (zero RPM under load)
- Configurable temperature thresholds
- Background service with Windows startup support

## 📥 Installation

Download the [latest release](https://github.com/kh4f/amdalert/releases/latest) and extract it anywhere.

## 🕹️ Usage

- Control the app via CLI (`AMDAlert.exe`).
- Configuration is stored in `settings.json` (created automatically on first run). You can edit it manually or via CLI commands.
- Enabling alerts starts a background service (`AMDAlertDaemon.exe`) that checks the GPU state every 10 sec and shows notifications when needed.

</br>

<div align="center">
  <b>MIT License © 2026 <a href="https://github.com/kh4f">kh4f</a></b>
</div>