<div align="center">
	<img alt="logo" src="assets/logo.png">
	<br>
	<a href="https://github.com/kh4f/amdalert/releases"><img src="https://img.shields.io/github/v/tag/kh4f/amdalert?label=%F0%9F%8F%B7%EF%B8%8F%20Release&style=flat-square&color=ceccff&labelColor=303145" alt="version"/></a>&nbsp;
	<a href="https://github.com/kh4f/amdalert/issues?q=is%3Aissue+is%3Aopen+label%3Abug"><img src="https://img.shields.io/github/issues/kh4f/amdalert/bug?label=%F0%9F%90%9B%20Bugs&style=flat-square&color=ceccff&labelColor=303145" alt="bugs"></a>&nbsp;
	<a href="https://github.com/kh4f/amdalert/blob/master/LICENSE"><img src="https://img.shields.io/github/license/kh4f/amdalert?style=flat-square&label=%F0%9F%9B%A1%EF%B8%8F%20License&color=ceccff&labelColor=303145" alt="license"></a>&nbsp;
	<br><br>
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
- Configuration is stored in `config.json` (created automatically on first run). You can edit it manually or via CLI commands.
- Enabling alerts starts a background service (`AMDAlertDaemon.exe`) that checks the GPU state every 10 sec and shows notifications when needed.

</br>

<div align="center">
  <b>MIT License © 2026 <a href="https://github.com/kh4f">kh4f</a></b>
</div>