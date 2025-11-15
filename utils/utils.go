package utils

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <android/log.h>
#cgo arm64 LDFLAGS: -L../libs/arm64-v8a -lAutoGo
#cgo amd64 LDFLAGS: -L../libs/x86_64 -lAutoGo
#cgo 386 LDFLAGS: -L../libs/x86 -lAutoGo

void logI(const char* label,const char* message) {
    __android_log_print(ANDROID_LOG_INFO, label, "%s", message);
}

void logE(const char* label, const char* message) {
    __android_log_print(ANDROID_LOG_ERROR, label, "%s", message);
}

char* shell(const char* cmd);
*/
import "C"
import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/xiaocainiao633/Genie1.0--/utils/java"
)

func LogI(label, message string) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.logI(cLabel, cMessage)
}

func LogE(label, message string) {
	cLabel := C.CString(label)
	defer C.free(unsafe.Pointer(cLabel))
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cMessage))
	C.logE(cLabel, cMessage)
}

// 执行shell命令
func Shell(cmd string) string {
	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))

	cResult := C.shell(cCmd)
	defer C.free(unsafe.Pointer(cResult))

	// 返回shell命令执行结果
	return C.GoString(cResult)
}

func Random(min, max int) int {
	if min > max {
		min, max = max, min // 确保 min <= max
	}

	delta := max - min + 1 // 计算范围（包含 max）
	if delta <= 0 {
		return min // 当 min == max 时直接返回
	}

	// 生成 0 到 delta-1 的随机数
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(delta)))
	if err != nil {
		panic(err) // 处理错误（例如系统随机源不可用）
	}

	return min + int(nBig.Int64()) // 转换为 min ~ max 的范围
}

func Sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}

func I2s(i int) string {
	return strconv.Itoa(i)
}

func S2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

// F2s 将浮点数转换为字符串。
func F2s(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// S2f 将字符串转换为浮点数。如果转换失败返回0.0。
func S2f(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

// B2s 将布尔值转换为字符串 ("true" 或 "false")。
func B2s(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// S2b 将字符串转换为布尔值。如果无法转换则返回 false。
func S2b(s string) bool {
	b, _ := strconv.ParseBool(strings.TrimSpace(s))
	return b
}

// CallJavaMethod 调用Java方法
func CallJavaMethod(method, params string) string {
	cMethod := C.CString(method)
	defer C.free(unsafe.Pointer(cMethod))
	cParams := C.CString(params)
	defer C.free(unsafe.Pointer(cParams))
	cResult := C.callJavaMethod(cMethod, cParams)
	defer C.free(unsafe.Pointer(cResult))
	return C.GoString(cResult)
}