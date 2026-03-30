// shooter 使用 SDL2_mixer 在树莓派等 Linux 环境播放射击相关音效（演示 CLI）。
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jammy/voice_play/internal/audio"
)

func main() {
	assets := flag.String("assets", "assets/sounds", "含 shot.wav / reload.wav / helmet.wav 的目录")
	demo := flag.Bool("demo", false, "运行一段演示后退出")
	autoInterval := flag.Duration("autofire-interval", 120*time.Millisecond, "演示中连发间隔")
	interactiveAuto := flag.Duration("interactive-autofire", 120*time.Millisecond, "交互模式下 a 命令使用的连发间隔")
	flag.Parse()

	cfg := audio.DefaultConfig(*assets)
	if err := audio.EnsureAssets(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "资源检查失败: %v\n", err)
		os.Exit(1)
	}

	eng, err := audio.NewEngine(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化音频: %v\n", err)
		os.Exit(1)
	}
	defer eng.Close()

	if *demo {
		runDemo(eng, *autoInterval)
		return
	}

	runInteractive(eng, *interactiveAuto)
}

func runDemo(eng *audio.Engine, autoInterval time.Duration) {
	fmt.Println("演示: 单发 -> 连发+换弹叠加 -> 头盔叠加")
	_ = eng.PlayShot()
	time.Sleep(200 * time.Millisecond)

	eng.StartAutoFire(autoInterval)
	time.Sleep(100 * time.Millisecond)
	_ = eng.PlayReload()
	time.Sleep(300 * time.Millisecond)
	_ = eng.PlayHelmet()
	time.Sleep(200 * time.Millisecond)
	_ = eng.PlayHelmet()

	time.Sleep(800 * time.Millisecond)
	eng.StopAutoFire()
	fmt.Println("演示结束")
}

func runInteractive(eng *audio.Engine, autoEvery time.Duration) {
	fmt.Printf("交互命令: s=射击 r=换弹 h=头盔 a=开始连发 x=停止连发 q=退出（连发间隔 %v）\n", autoEvery)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sc := bufio.NewScanner(os.Stdin)
	lines := make(chan string)
	go func() {
		for sc.Scan() {
			lines <- strings.TrimSpace(sc.Text())
		}
	}()

	for {
		select {
		case <-sigCh:
			eng.StopAutoFire()
			return
		case line := <-lines:
			if line == "" {
				continue
			}
			switch line {
			case "q", "quit":
				eng.StopAutoFire()
				return
			case "s":
				if err := eng.PlayShot(); err != nil {
					fmt.Fprintln(os.Stderr, "射击:", err)
				}
			case "r":
				if err := eng.PlayReload(); err != nil {
					fmt.Fprintln(os.Stderr, "换弹:", err)
				}
			case "h":
				if err := eng.PlayHelmet(); err != nil {
					fmt.Fprintln(os.Stderr, "头盔:", err)
				}
			case "a":
				eng.StartAutoFire(autoEvery)
				fmt.Println("连发已开始")
			case "x":
				eng.StopAutoFire()
				fmt.Println("连发已停止")
			default:
				fmt.Println("未知命令:", line)
			}
		}
	}
}
