package console

/*
#include "imgui.h"
#include <stdlib.h>
#cgo arm64 LDFLAGS: -L../../resources/libs/arm64-v8a -limgui
#cgo amd64 LDFLAGS: -L../../resources/libs/x86_64 -limgui
#cgo 386 LDFLAGS: -L../../resources/libs/x86 -limgui
*/
import "C"
import (
	"fmt"
	"github.com/xiaocainiao633/Genie1.0--/device"
	"github.com/xiaocainiao633/Genie1.0--/https"
	"github.com/xiaocainiao633/Genie1.0--/utils"
	"os"
	"runtime"
	"strings"
	"unsafe"
)

var isInit = false

// 初始化GUI系统
func Init(noCaptureMode bool) {
	isInit = true
	if int(C.CheckInit()) == 0 {
		if runtime.GOARCH != "arm64" {
			noCaptureMode = true
			_, err := os.Stat("/system/fonts/NotoSansCJK-Regular.ttc")
			if err != nil {
				_, err = os.Stat("/system/fonts/NotoSerifCJK-Regular.ttc")
				if err != nil {
					_, err = os.Stat("/data/local/tmp/NotoSansCJK-Regular.ttc")
					if err != nil {
						fmt.Println("imgui初始化中..")
						code, data := https.Get("https://vip.123pan.cn/1823847070/AutoGo/NotoSansCJK-Regular.ttc", 20000)
						if code == 200 {
							fmt.Println("imgui初始化完毕")
							os.WriteFile("/data/local/tmp/NotoSansCJK-Regular.ttc", data, 0644)
						} else {
							fmt.Println("imgui初始失败,下载依赖文件超时,中文可能显示乱码")
						}
					}
				}
			}
		}
		go C.Init(C.int(b2i(noCaptureMode)))
		success := false
		for i := 0; i < 100; i++ {
			utils.Sleep(10)
			if int(C.CheckInit()) == 1 {
				success = true
				break
			}
		}
		if !success {
			fmt.Fprintf(os.Stderr, "[AutoGo] imgui初始化失败")
			os.Exit(1)
		}
		w := intMin(device.Width, device.Height)
		h := intMax(device.Width, device.Height)
		if w < 1080 {
			C.Toast_setTextSize(35)
		} else if h > 1920 {
			C.Toast_setTextSize(50)
		}
		scale := float32(1080) / float32(w)
		C.Console_setPosition(25, C.int(int(float32(50)/scale)))
		if device.Width > device.Height {
			C.Console_setSize(C.int(int(float32(800)/scale)), C.int(int(float32(520)/scale)))
		} else {
			C.Console_setSize(C.int(int(float32(520)/scale)), C.int(int(float32(800)/scale)))
		}
	}
}

// 设置控制台宽高
func SetWindowSize(width, height int) {
	C.Console_setSize(C.int(width), C.int(height))
}

// 设置控制台位置
func SetWindowPosition(x, y int) {
	C.Console_setPosition(C.int(x), C.int(y))
}

// 设置控制台背景色 
func SetWindowColor(color string) {
	cColor := C.CString(color)
	defer C.free(unsafe.Pointer(cColor))
	C.Console_setWindowColor(cColor)
}

// 设置文字颜色
func SetTextColor(color string) {
	cColor := C.CString(color)
	defer C.free(unsafe.Pointer(cColor))
	C.Console_setTextColor(cColor)
}

// SetTextSize 设置字体大小
func SetTextSize(size int) {
	C.Console_setTextSize(C.int(size))
}

// 输出文本到控制台
func Println(a ...any) {
	str := fmt.Sprint(a...)
	arr := strings.Split(str, "\n")
	for _, line := range arr {
		cLine := C.CString(line)
		defer C.free(unsafe.Pointer(cLine))
		C.Console_println(cLine)
	}
}

// 清空控制台
func Clear() {
	C.Console_clear()
}

// 显示控制台
func Show() {
	if !isInit {
		Init(false)
	}
	C.Console_show()
}

// 隐藏控制台
func Hide() {
	C.Console_hide()
}

// 布尔值转整数
func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// 最小值
func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 最大值
func intMax(a, b int) int {
	if a < b {
		return b
	}
	return a
}
