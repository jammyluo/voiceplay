#!/usr/bin/env bash
# 在项目根目录构建 shooter 二进制（需已安装 libsdl2 / libsdl2-mixer 与 Go）。
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
export CGO_ENABLED=1
go build -o bin/shooter ./cmd/shooter
echo "输出: $ROOT/bin/shooter"
