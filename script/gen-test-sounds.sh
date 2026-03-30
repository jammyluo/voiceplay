#!/usr/bin/env bash
# 用 ffmpeg 生成占位 WAV（便于无素材时跑通流程）。需已安装 ffmpeg。
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="${1:-$ROOT/assets/sounds}"
mkdir -p "$OUT"

if ! command -v ffmpeg >/dev/null 2>&1; then
  echo "请先安装 ffmpeg: sudo apt install ffmpeg 或 brew install ffmpeg" >&2
  exit 1
fi

# 短促噪声/音调，便于区分：射击=短脉冲，换弹=较长滑音，头盔=中频叮
ffmpeg -y -f lavfi -i "sine=frequency=800:duration=0.05" -ac 2 -ar 44100 "$OUT/shot.wav" >/dev/null 2>&1
ffmpeg -y -f lavfi -i "sine=frequency=200:duration=0.6" -ac 2 -ar 44100 "$OUT/reload.wav" >/dev/null 2>&1
ffmpeg -y -f lavfi -i "sine=frequency=1200:duration=0.15" -ac 2 -ar 44100 "$OUT/helmet.wav" >/dev/null 2>&1

echo "已写入: $OUT/shot.wav, $OUT/reload.wav, $OUT/helmet.wav"
