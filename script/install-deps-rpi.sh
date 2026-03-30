#!/usr/bin/env bash
# 在 Raspberry Pi OS / Debian 上安装 SDL2 + SDL2_mixer 开发包（构建 CGO 所需）。
set -euo pipefail
sudo apt-get update
sudo apt-get install -y \
  build-essential \
  pkg-config \
  libsdl2-dev \
  libsdl2-mixer-dev
