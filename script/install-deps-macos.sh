#!/usr/bin/env bash
# macOS 上通过 Homebrew 安装 SDL2 / SDL2_mixer（用于本机 CGO 编译与调试）。
set -euo pipefail
brew install sdl2 sdl2_mixer pkg-config
