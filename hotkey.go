package main

import (
	"syscall"
	"unsafe"
)

// var (
// 	registerHotKey    *windows.LazyProc
// 	unregisterHotKey  *windows.LazyProc
// 	postThreadMessage *windows.LazyProc

// 	getCurrentThread *windows.LazyProc
// 	getThreadId      *windows.LazyProc
// )

// func init() {
// 	// Library
// 	libuser32 := windows.NewLazySystemDLL("user32.dll")
// 	libkernel32 := windows.NewLazySystemDLL("kernel32.dll")

// 	// Functions
// 	registerHotKey = libuser32.NewProc("RegisterHotKey")
// 	unregisterHotKey = libuser32.NewProc("UnregisterHotKey")
// 	postThreadMessage = libuser32.NewProc("PostThreadMessageW")

// 	getCurrentThread = libkernel32.NewProc("GetCurrentThread")
// 	getThreadId = libkernel32.NewProc("GetThreadId")
// }

var (
	moduser32            = syscall.NewLazyDLL("user32.dll")
	procGetMessage       = moduser32.NewProc("GetMessageW")
	procRegisterHotKey   = moduser32.NewProc("RegisterHotKey")
	procUnregisterHotKey = moduser32.NewProc("UnregisterHotKey")
)

const (
	MOD_ALT      = 0x0001
	MOD_CONTROL  = 0x0002
	MOD_NOREPEAT = 0x4000

	// https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
	VK_P = 0x50
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

func GetMessage(msg *MSG, hwnd uintptr, msgFilterMin, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))

	return int(ret)
}

func RegisterHotKey(id int, fsModifiers uint, vk uint) bool {
	ret, _, _ := procRegisterHotKey.Call(
		0,
		uintptr(id),
		uintptr(fsModifiers),
		uintptr(vk),
	)
	return ret != 0
}

func UnregisterHotKey(id int32) bool {
	ret, _, _ := procUnregisterHotKey.Call(
		0,
		uintptr(id))

	return ret != 0
}

// https://github.com/willemvds/Hopp-polla/blob/ac44bbbbf5ac772b60e3265e3f48c6dc535b289c/eventserver/events_windows.go
func StartHotkeyListener(done chan bool) {
	go func() {
		var msg MSG
		RegisterHotKey(1, MOD_ALT|MOD_CONTROL|MOD_NOREPEAT, VK_P)
		for {
			ok := GetMessage(&msg, 0, 0, 0)
			if ok != 1 {
				continue
			}
			switch msg.WParam {
			case 1:
				done <- true
			default:

			}
		}
	}()

}
