package hotkey

import (
	"fmt"
	"sync"
	"time"

	"github.com/robotn/gohook"
)

// HotkeyManager 管理全局快捷键的类型
type HotkeyManager struct {
	running  bool
	mu       sync.Mutex
	toggleCb func() // 切换窗口回调函数
}

// NewHotkeyManager 创建一个新的快捷键管理器
func NewHotkeyManager(toggleCb func()) *HotkeyManager {
	return &HotkeyManager{
		running:  false,
		toggleCb: toggleCb,
	}
}

// Start 注册并启动全局快捷键监听
func (hm *HotkeyManager) Start() error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	if hm.running {
		return fmt.Errorf("热键管理器已经在运行中")
	}

	clickTime := time.Now()
	hook.Register(hook.KeyDown, []string{"cmd", "c"}, func(e hook.Event) {
		if time.Since(clickTime) < time.Millisecond*300 {
			if hm.toggleCb != nil {
				hm.toggleCb()
			}
		}
		clickTime = time.Now()
	})

	// 启动钩子监听
	evChan := hook.Start()
	go func() {
		<-hook.Process(evChan)
	}()

	hm.running = true
	return nil
}

// Stop 停止全局快捷键监听
func (hm *HotkeyManager) Stop() {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	if hm.running {
		hook.End()
		hm.running = false
		fmt.Println("已停止全局快捷键监听")
	}
}
