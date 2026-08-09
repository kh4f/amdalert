<div align="center">
	<picture>
		<source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/kh4f/amdlert/refs/heads/assets/logo-dark.png">
		<img alt="logo" src="https://raw.githubusercontent.com/kh4f/amdlert/refs/heads/assets/logo-light.png">
	</picture>
	<br>
	A lightweight <b>watchdog</b> for your <b>AMD GPU temperature</b>
	<br><br>
	<p>
		<a href='https://github.com/kh4f/amdlert/releases'><img alt="release" src="https://img.shields.io/github/v/tag/kh4f/amdlert?style=flat-square&labelColor=FC0042&color=D5CAFD&label=%F0%9F%8F%B7%EF%B8%8F%20release"></a>&nbsp;
		<a href="https://github.com/kh4f/amdlert/releases"><img alt="downloads" src="https://img.shields.io/github/downloads/kh4f/amdlert/total?style=flat-square&labelColor=FC0042&color=D5CAFD&label=%F0%9F%93%A5%20downloads" /></a>&nbsp;
		<a href="https://github.com/kh4f/amdlert/blob/main/LICENSE"><img alt="license" src="https://img.shields.io/github/license/kh4f/amdlert?style=flat-square&labelColor=FC0042&color=D5CAFD&label=%F0%9F%9B%A1%EF%B8%8F%20license"></a>
	</p>
	<b>
		<a href="#-overview">Overview</a>&nbsp; •&nbsp;
		<a href="#-install">Install</a>&nbsp; •&nbsp;
		<a href="#%EF%B8%8F-usage">Usage</a>
	</b>
	<br><br>
    <img alt="demo" src="https://raw.githubusercontent.com/kh4f/amdlert/refs/heads/assets/demo.png">
</div>

## 👀 Overview

**AMDlert** is a single-binary CLI + daemon for Windows that monitors your AMD GPU temperature every 10s and alerts you when it exceeds a configurable threshold.

## 📥 Install

1. Download and extract the [latest release](https://github.com/kh4f/amdlert/releases) `.zip`
2. Move the `AMDlert` folder to your preferred location (e.g., `C:/Program Files/`)
3. Create a shortcut for `AMDlert.exe` or pin it to the Start menu for faster access

## 🕹️ Usage

Run `AMDlert.exe` to open the interactive menu:
- **Run daemon** — start background monitoring (checks every 10s)
- **Set threshold** — change the temperature (°C) that triggers alerts
- **Enable autostart** — launch the daemon automatically on Windows startup

The threshold is stored in `config.json` next to the executable (created automatically; default is `50`).