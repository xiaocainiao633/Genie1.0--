package motion

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"github.com/Dasongzi1366/AutoGo/utils"
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

// TouchDown 触摸按下
func TouchDown(x, y, fingerID int) { // 参数：x, y, fingerID
	// 确保 fingerID 在有效范围内
	fingerID = fingerID - 1
	if fingerID < 0 || fingerID > 9 {
		fingerID = 0
	}
	// 发送触摸按下事件
	utils.Send(fmt.Sprintf("d|%d|%d|%d", x, y, fingerID))
}

// TouchMove 触摸移动
func TouchMove(x, y, fingerID int) {
	fingerID = fingerID - 1
	if fingerID < 0 || fingerID > 9 {
		fingerID = 0
	}
	utils.Send(fmt.Sprintf("m|%d|%d|%d", x, y, fingerID))
}

// TouchUp 触摸抬起
func TouchUp(x, y, fingerID int) {
	fingerID = fingerID - 1
	if fingerID < 0 || fingerID > 9 {
		fingerID = 0
	}
	utils.Send(fmt.Sprintf("u|%d|%d|%d", x, y, fingerID))
}

// Click 点击
func Click(x, y, fingerID int) {
	TouchDown(x, y, fingerID)
	sleep(random(10, 20))
	TouchUp(x, y, fingerID)
}

// LongClick 长按
func LongClick(x, y, duration int) {
	TouchDown(x, y, 1)
	sleep(duration + random(1, 20))
	TouchUp(x, y, 1)
}

// Swipe 滑动
func Swipe(x1, y1, x2, y2, duration int) {
	utils.Send(fmt.Sprintf("s1|%d|%d|%d|%d|%d", x1, y1, x2, y2, duration))
}

// Swipe2 滑动
func Swipe2(x1, y1, x2, y2, duration int) {
	utils.Send(fmt.Sprintf("s2|%d|%d|%d|%d|%d", x1, y1, x2, y2, duration))
}

// Home 点击Home键
func Home() {
	KeyAction(KEYCODE_HOME)
}

// Back 点击Back键
func Back() {
	KeyAction(KEYCODE_BACK)
}

// Recents 点击应用切换键
func Recents() {
	KeyAction(KEYCODE_APP_SWITCH)
}

// PowerDialog 长按电源键
func PowerDialog() {
	shell("input keyevent --longpress KEYCODE_POWER")
}

// Notifications 点击通知键
func Notifications() {
	KeyAction(KEYCODE_NOTIFICATION)
}

// QuickSettings 点击快速设置键
func QuickSettings() {
	shell("cmd statusbar expand-settings")
}

// VolumeUp 点击音量增加键
func VolumeUp() {
	KeyAction(KEYCODE_VOLUME_UP)
}

// VolumeDown 点击音量减小键
func VolumeDown() {
	KeyAction(KEYCODE_VOLUME_DOWN)
}

// Camera 点击相机键
func Camera() {
	KeyAction(KEYCODE_CAMERA)
}

// KeyAction 模拟按键
func KeyAction(code int) {
	// 发送按键事件
	utils.Send(fmt.Sprintf("k|%d", code))
}

// random 生成随机数
func random(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

// sleep 睡眠
func sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}

// shell 执行shell命令
func shell(cmd string) {
	if strings.Contains(cmd, ";") {
		cmd = "(" + cmd + ")"
	}
	cCmd := C.CString(cmd + " > /dev/null 2>&1")
	defer C.free(unsafe.Pointer(cCmd))
	C.system(cCmd)
}
