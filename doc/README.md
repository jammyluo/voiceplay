# voice_play

树莓派等 Linux 环境下，使用 **Go + SDL2_mixer** 播放射击类音效：单发、连发、换弹、命中头盔；头盔可与其它音效**同时播放**（mixer 多通道混音）。

## 依赖

- Go 1.21+
- CGO：`libsdl2`、`libsdl2-mixer`
  - **树莓派 / Debian**：`script/install-deps-rpi.sh`
  - **macOS（本机试编译）**：`brew install sdl2 sdl2_mixer pkg-config`
- 运行目录下需有 `assets/sounds/shot.wav`、`reload.wav`、`helmet.wav`（可用 `script/gen-test-sounds.sh` 生成占位文件）

## 构建

```bash
chmod +x script/build.sh
./script/build.sh
```

## 运行

```bash
./bin/shooter -demo
./bin/shooter
```

交互模式：`s` 射击，`r` 换弹，`h` 头盔，`a` 连发，`x` 停止连发，`q` 退出。

更多说明见 [architecture.md](./architecture.md)。
