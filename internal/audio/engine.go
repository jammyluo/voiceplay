// Package audio 封装 SDL2_mixer：多路 Chunk 混音，支持射击、连发、换弹、头盔叠加。
package audio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	defaultFrequency = 44100
	defaultChunkSize = 2048
	defaultChannels  = 32
)

// Engine 持有预加载的音效；Play* 可并发调用，内部串行交给 mixer。
type Engine struct {
	shot   *mix.Chunk
	reload *mix.Chunk
	helmet *mix.Chunk

	mu         sync.Mutex
	autoCancel context.CancelFunc
	autoDone   sync.WaitGroup
}

// Config 资源文件名相对于 AssetDir。
type Config struct {
	AssetDir   string
	ShotFile   string
	ReloadFile string
	HelmetFile string
}

// DefaultConfig 使用约定文件名。
func DefaultConfig(assetDir string) Config {
	return Config{
		AssetDir:   assetDir,
		ShotFile:   "ak-47-single-shot_1.wav",
		ReloadFile: "ak-47-reloaded.wav",
		HelmetFile: "helmet.wav",
	}
}

// NewEngine 初始化 SDL 音频、打开 mixer、加载 WAV。须在同一进程内 Close。
func NewEngine(cfg Config) (*Engine, error) {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		return nil, fmt.Errorf("sdl.Init: %w", err)
	}
	if err := mix.Init(0); err != nil {
		sdl.Quit()
		return nil, fmt.Errorf("mix.Init: %w", err)
	}
	if err := mix.OpenAudio(defaultFrequency, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, defaultChunkSize); err != nil {
		mix.Quit()
		sdl.Quit()
		return nil, fmt.Errorf("mix.OpenAudio: %w", err)
	}
	_ = mix.AllocateChannels(defaultChannels)

	e := &Engine{}
	var err error
	path := func(name string) string {
		return filepath.Join(cfg.AssetDir, name)
	}
	if e.shot, err = mix.LoadWAV(path(cfg.ShotFile)); err != nil {
		e.closeMixer()
		return nil, fmt.Errorf("load %s: %w", cfg.ShotFile, err)
	}
	if e.reload, err = mix.LoadWAV(path(cfg.ReloadFile)); err != nil {
		e.closeMixer()
		return nil, fmt.Errorf("load %s: %w", cfg.ReloadFile, err)
	}
	if e.helmet, err = mix.LoadWAV(path(cfg.HelmetFile)); err != nil {
		e.closeMixer()
		return nil, fmt.Errorf("load %s: %w", cfg.HelmetFile, err)
	}
	return e, nil
}

func (e *Engine) closeMixer() {
	if e.helmet != nil {
		e.helmet.Free()
		e.helmet = nil
	}
	if e.reload != nil {
		e.reload.Free()
		e.reload = nil
	}
	if e.shot != nil {
		e.shot.Free()
		e.shot = nil
	}
	mix.CloseAudio()
	mix.Quit()
	sdl.Quit()
}

// Close 停止连发并释放资源。
func (e *Engine) Close() {
	e.StopAutoFire()
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.helmet != nil {
		e.helmet.Free()
		e.helmet = nil
	}
	if e.reload != nil {
		e.reload.Free()
		e.reload = nil
	}
	if e.shot != nil {
		e.shot.Free()
		e.shot = nil
	}
	mix.CloseAudio()
	mix.Quit()
	sdl.Quit()
}

func (e *Engine) playChunk(ch *mix.Chunk) error {
	if ch == nil {
		return nil
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, err := ch.Play(-1, 0); err != nil {
		return err
	}
	return nil
}

// PlayShot 播放单发；可与换弹、头盔叠加。
func (e *Engine) PlayShot() error {
	return e.playChunk(e.shot)
}

// PlayReload 换弹。
func (e *Engine) PlayReload() error {
	return e.playChunk(e.reload)
}

// PlayHelmet 命中头盔；可与射击、换弹同时播放。
func (e *Engine) PlayHelmet() error {
	return e.playChunk(e.helmet)
}

// StartAutoFire 以 interval 间隔重复播放射击声，直至 StopAutoFire。再次调用会先停止上一轮。
func (e *Engine) StartAutoFire(interval time.Duration) {
	e.StopAutoFire()
	if interval <= 0 {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.autoCancel = cancel
	e.autoDone.Add(1)
	go func() {
		defer e.autoDone.Done()
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				_ = e.PlayShot()
			}
		}
	}()
}

// StopAutoFire 停止连发。
func (e *Engine) StopAutoFire() {
	if e.autoCancel != nil {
		e.autoCancel()
		e.autoDone.Wait()
		e.autoCancel = nil
	}
}

// EnsureAssets 检查约定 WAV 是否存在，便于启动前报错。
func EnsureAssets(cfg Config) error {
	names := []string{cfg.ShotFile, cfg.ReloadFile, cfg.HelmetFile}
	for _, n := range names {
		p := filepath.Join(cfg.AssetDir, n)
		st, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("音效文件: %w", err)
		}
		if st.IsDir() {
			return fmt.Errorf("音效路径是目录: %s", p)
		}
	}
	return nil
}
