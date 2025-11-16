package imgui

/*
#include "imgui.h"
#include <stdlib.h>
#cgo arm64 LDFLAGS: -L../../resources/libs/arm64-v8a -limgui
#cgo amd64 LDFLAGS: -L../../resources/libs/x86_64 -limgui
#cgo 386 LDFLAGS: -L../../resources/libs/x86 -limgui
*/
import "C"
import (
	"encoding/base64"
	"fmt"
	"github.com/xiaocainiao633/Genie1.0--/device"
	"github.com/xiaocainiao633/Genie1.0--/https"
	"github.com/xiaocainiao633/Genie1.0--/utils"
	"os"
	"runtime"
	"time"
	"unsafe"
)

type TextItem struct {
	TextColor string //文字颜色 例如 #FFFFFF
	Text      string
}

var isInit = false

// 初始化GUI
func Init(noCaptureMode bool) {
	if !isInit {
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
				utils.Sleep(20)
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
}

// 显示文字消息
func Toast(message string) {
	Init(false)
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.Toast(cMessage)
}

// 设置文字大小
func ToastSetTextSize(size int) {
	C.Toast_setTextSize(C.int(size))
}

// HUD初始化
func HudInit(x1, y1, x2, y2 int, bgColor string, textSize int) {
	Init(false)
	if textSize <= 0 {
		textSize = 45
	}
	cBgColor := C.CString(bgColor)
	defer C.free(unsafe.Pointer(cBgColor))
	C.Hud_init(C.int(x1), C.int(y1), C.int(x2), C.int(y2), cBgColor, C.int(textSize))
}

// 设置文本
func HudSetText(items []TextItem) {
	str := ""
	for _, item := range items {
		if item.Text != "" {
			str = str + item.TextColor + "|" + item.Text + "|"
		}
	}
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	C.Hud_setText(cStr)
}

// 关闭HUD
func HudClose() {
	C.Hud_clear()
}

// 图形绘制
func DrawRect(x1, y1, x2, y2 int, color string) {
	Init(false)
	cColor := C.CString(color)
	defer C.free(unsafe.Pointer(cColor))
	C.Rect_add(C.int(x1), C.int(y1), C.int(x2), C.int(y2), cColor)
}

func ClearRect() {
	C.Rect_clear()
}

// 直线绘制
func DrawLine(x1, y1, x2, y2 int, color string) {
	Init(false)
	cColor := C.CString(color)
	defer C.free(unsafe.Pointer(cColor))
	C.StrLine_add(C.int(x1), C.int(y1), C.int(x2), C.int(y2), cColor)
}

func ClearLine() {
	C.StrLine_clear()
}

// 图像显示
func DrawImg(x1, y1, x2, y2 int, pngData []byte) {
	Init(false)
	C.Image_add(
		C.int(x1),
		C.int(y1),
		C.int(x2),
		C.int(y2),
		(*C.uchar)(unsafe.Pointer(&pngData[0])),
		C.int(len(pngData)),
	)
}

func ClearImg() {
	C.Image_clear()
}

// 关闭GUI系统
func Close() {
	C.Close()
}

// 显示警告提示框
func Alert(title, message string) {
	if os.Getenv("APPPID") != "" {
		cCmd := C.CString("am broadcast -a com.autogo --es message alert --es title " + base64.StdEncoding.EncodeToString([]byte(title)) + " --es msg " + base64.StdEncoding.EncodeToString([]byte(message)))
		defer C.free(unsafe.Pointer(cCmd))
		C.system(cCmd)
		time.Sleep(1 * time.Second)
	} else {
		/*cTitle := C.CString(title)
		defer C.free(unsafe.Pointer(cTitle))
		cMessage := C.CString(message)
		defer C.free(unsafe.Pointer(cMessage))
		C.Alert(cTitle, cMessage)*/
	}
}

// 布尔转整数
func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// 求最小值
func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 求最大值
func intMax(a, b int) int {
	if a < b {
		return b
	}
	return a
}
