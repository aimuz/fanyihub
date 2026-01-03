package hotkey

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices -framework CoreFoundation -framework Foundation
#import <ApplicationServices/ApplicationServices.h>
#import <CoreFoundation/CoreFoundation.h>
#import <Foundation/Foundation.h>

bool checkAccessibility(bool prompt) {
    NSDictionary *opts = @{(__bridge NSString *)kAXTrustedCheckOptionPrompt: @(prompt)};
    return AXIsProcessTrustedWithOptions((__bridge CFDictionaryRef)opts);
}
*/
import "C"

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
)

// HotkeyManager 管理全局快捷键的类型
type HotkeyManager struct {
	running     bool
	mu          sync.Mutex
	toggleCb    func()        // 切换窗口回调函数
	ocrCb       func()        // OCR 截图回调函数
	statusCb    func(bool)    // 权限状态回调函数
	stopPolling chan struct{} // 停止轮询信号
	clickTime   time.Time     // 上次点击时间
}

// NewHotkeyManager 创建一个新的快捷键管理器
func NewHotkeyManager(toggleCb func(), ocrCb func()) *HotkeyManager {
	return &HotkeyManager{
		running:   false,
		toggleCb:  toggleCb,
		ocrCb:     ocrCb,
		clickTime: time.Now(),
	}
}

// SetStatusCallback 设置权限状态变更回调
func (hm *HotkeyManager) SetStatusCallback(cb func(bool)) {
	hm.statusCb = cb
}

// IsAccessibilityEnabled 检查辅助功能权限是否已授予
// prompt: 是否弹出系统授权提示
func IsAccessibilityEnabled(prompt bool) bool {
	return bool(C.checkAccessibility(C.bool(prompt)))
}

// Start 注册并启动全局快捷键监听
func (hm *HotkeyManager) Start() error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	if hm.running {
		return fmt.Errorf("热键管理器已经在运行中")
	}

	// 检查辅助功能权限，如果没有则弹出提示
	if !IsAccessibilityEnabled(true) {
		slog.Warn("辅助功能权限未授予，等待用户授权")
		// 通知前端权限状态
		if hm.statusCb != nil {
			hm.statusCb(false)
		}
		// 启动权限轮询
		hm.startPermissionPolling()
		return nil
	}

	// 权限已授予，直接启动
	return hm.startHook()
}

// startHook 内部方法：启动 hook 监听
func (hm *HotkeyManager) startHook() error {
	hook.Register(hook.KeyDown, []string{"cmd", "c"}, func(e hook.Event) {
		if time.Since(hm.clickTime) < time.Millisecond*300 {
			if hm.toggleCb != nil {
				hm.toggleCb()
			}
		}
		hm.clickTime = time.Now()
	})

	// 注册 OCR 截图快捷键: Cmd+Shift+O
	hook.Register(hook.KeyDown, []string{"cmd", "shift", "o"}, func(e hook.Event) {
		if hm.ocrCb != nil {
			hm.ocrCb()
		}
	})

	// 启动钩子监听
	evChan := hook.Start()
	go func() {
		<-hook.Process(evChan)
	}()

	hm.running = true
	slog.Info("全局快捷键已启动")

	// 通知前端权限已授予
	if hm.statusCb != nil {
		hm.statusCb(true)
	}

	return nil
}

// startPermissionPolling 启动权限轮询
func (hm *HotkeyManager) startPermissionPolling() {
	hm.stopPolling = make(chan struct{})

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		slog.Info("开始轮询辅助功能权限状态")

		for {
			select {
			case <-hm.stopPolling:
				slog.Info("停止权限轮询")
				return
			case <-ticker.C:
				// 检查权限（不弹出提示）
				if IsAccessibilityEnabled(false) {
					slog.Info("检测到辅助功能权限已授予，正在启动快捷键监听")

					hm.mu.Lock()
					if !hm.running {
						if err := hm.startHook(); err != nil {
							slog.Error("启动快捷键监听失败", "error", err)
						}
					}
					hm.mu.Unlock()

					// 权限已授予，停止轮询
					return
				}
			}
		}
	}()
}

// Stop 停止全局快捷键监听
func (hm *HotkeyManager) Stop() {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	// 停止轮询
	if hm.stopPolling != nil {
		close(hm.stopPolling)
		hm.stopPolling = nil
	}

	if hm.running {
		hook.End()
		hm.running = false
		slog.Info("已停止全局快捷键监听")
	}
}
