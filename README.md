<div align="center">
	<picture>
		<source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/kh4f/amdalert/refs/heads/assets/logo-dark.png">
		<img alt="logo" src="https://raw.githubusercontent.com/kh4f/amdalert/refs/heads/assets/logo-light.png">
	</picture>
	<br>
	A lightweight <b>watchdog</b> for your <b>AMD GPU temperature</b>
	<br><br>
	<p>
		<a href="https://github.com/kh4f/amdalert/releases"><img src="https://img.shields.io/github/v/tag/kh4f/amdalert?label=%F0%9F%8F%B7%EF%B8%8F%20Release&style=flat-square&color=D5CAFC&labelColor=B60C36" alt="version"/></a>&nbsp;
		<a href="https://github.com/kh4f/amdalert/issues?q=is%3Aissue+is%3Aopen+label%3Abug"><img src="https://img.shields.io/github/issues/kh4f/amdalert/bug?label=%F0%9F%90%9B%20Bugs&style=flat-square&color=D5CAFC&labelColor=B60C36" alt="bugs"></a>&nbsp;
		<a href="https://github.com/kh4f/amdalert/blob/main/LICENSE"><img src="https://img.shields.io/github/license/kh4f/amdalert?style=flat-square&label=%F0%9F%9B%A1%EF%B8%8F%20License&color=D5CAFC&labelColor=B60C36" alt="license"></a>
	</p>
	<p><b>
		<a href="#-overview">Overview</a>&nbsp; •&nbsp;
		<a href="#-installation">Installation</a>&nbsp; •&nbsp;
		<a href="#%EF%B8%8F-usage">Usage</a>
	</b></p>
	<img alt="demo" src="https://raw.githubusercontent.com/kh4f/amdalert/refs/heads/assets/demo.png">
</div>

## 👀 Overview
**AMDAlert** is a Windows utility that monitors your AMD GPU and warns you about overheating.

It’s especially useful for older cards, where the fan may silently stop.

### 🔥 Features
- Monitors GPU state and alerts on overheating
- Detects fan failures (e.g. zero RPM under load)
- Simple YAML-based configuration
- Launch at Windows startup

## 📥 Installation
Download and extract the [latest release](https://github.com/kh4f/amdalert/releases/latest).

## 🕹️ Usage
- Control the app from the command line with `AMDAlert.exe`.
- Configuration is stored in `config.yml` (created automatically on first run). You can edit it manually or via CLI сommands.
- Enabling alerts launches `AMDAlertDaemon.exe` in the background, that runs every 10 seconds, reloading the config if updated and showing notifications when needed.