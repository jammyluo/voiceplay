# 架构说明（SDL2_mixer）

## 分层

| 层 | 职责 |
|----|------|
| `cmd/shooter` | CLI：演示、交互、信号退出 |
| `internal/audio` | 初始化 SDL 音频、`mix.OpenAudio`、预加载 `mix.Chunk`、对外 `Play*` / 连发调度 |

## 混音与并发

- `mix.AllocateChannels(32)` 提供足够通道，多路 `Chunk.Play(-1, 0)` 由 SDL2_mixer 混音。
- 头盔、换弹、射击互不抢占「唯一声道」；若通道耗尽，`Play` 会返回错误（可调大通道数）。
- `Engine` 内 `sync.Mutex` 保护所有 `Play`，避免与 SDL 文档中线程约定冲突。

## 连发

- `StartAutoFire(interval)` 使用 `time.Ticker` + `PlayShot()`，与换弹/头盔调用路径相同，可自然叠加。
- `StopAutoFire` 取消 context 并 `Wait` 结束协程。

## 资源

- 使用 **WAV**（见 `mix.LoadWAV`），在 `OpenAudio` 之后加载，便于加载时完成格式转换。
- 文件名由 `audio.Config` 配置，默认 `shot.wav` / `reload.wav` / `helmet.wav`。

## 生命周期

1. `sdl.Init(INIT_AUDIO)` → `mix.Init(0)` → `mix.OpenAudio` → `LoadWAV` ×3  
2. 退出：`StopAutoFire` → `Chunk.Free` → `mix.CloseAudio` → `mix.Quit` → `sdl.Quit`

## 树莓派注意

- 无桌面环境时需确认 ALSA/Pulse 设备可用；若无声，检查 `aplay -l` 与默认输出。
- 交叉编译 ARM 时需匹配目标机的 SDL2 与 C 工具链；最简方式是在派上直接 `go build`。
