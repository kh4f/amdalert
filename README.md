<div align="center">
	<img alt="logo" src="assets/logo.png">
	<br>
	<a href="https://github.com/kh4f/amdalert/releases"><img src="https://img.shields.io/github/v/tag/kh4f/amdalert?label=%F0%9F%8F%B7%EF%B8%8F%20Release&style=flat-square&color=CECCFF&labelColor=303145" alt="version"/></a>&nbsp;
	<a href="https://github.com/kh4f/amdalert/issues?q=is%3Aissue+is%3Aopen+label%3Abug"><img src="https://img.shields.io/github/issues/kh4f/amdalert/bug?label=%F0%9F%90%9B%20Bugs&style=flat-square&color=CECCFF&labelColor=303145" alt="bugs"></a>&nbsp;
	<a href="https://github.com/kh4f/amdalert/blob/master/LICENSE"><img src="https://img.shields.io/github/license/kh4f/amdalert?style=flat-square&label=%F0%9F%9B%A1%EF%B8%8F%20License&color=CECCFF&labelColor=303145" alt="license"></a>&nbsp;
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
- Simple YAML configuration
- Launch at Windows startup

## 📥 Installation

Download and extract the [latest release](https://github.com/kh4f/amdalert/releases/latest).

## 🕹️ Usage

- Control the app from the command line with `AMDAlert.exe`.
- Configuration is stored in `config.yml` (created automatically on first run). You can edit it manually or via CLI сommands.
- Enabling alerts launches `AMDAlertDaemon.exe` in the background, that runs every 10 sec, reloading the config if updated and showing notifications when needed.