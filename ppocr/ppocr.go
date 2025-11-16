package ppocr

/*
#include "ppocr.h"
#include <stdlib.h>
#cgo arm64 LDFLAGS: -L../../resources/libs/arm64-v8a -lopencv_core -lopencv_imgproc -lppocr
#cgo amd64 LDFLAGS: -L../../resources/libs/x86_64 -lopencv_core -lopencv_imgproc -lppocr
#cgo 386 LDFLAGS: -L../../resources/libs/x86 -lopencv_core -lopencv_imgproc -lppocr
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/xiaocainiao633/Genie1.0--/files"
	"image"
	"os"
	"strings"
	"sync"
	"unsafe"

	"github.com/xiaocainiao633/Genie1.0--/images"
)

// 通过 CGO 调用 C/C++ 库实现文字识别

var pointer *C.Ppocr
var mutex sync.Mutex


// 识别结果包括位置坐标，文字内容，识别精度，中心点
type Result struct {
	X       int     `json:"X"`
	Y       int     `json:"Y"`
	Width   int     `json:"宽"`
	Height  int     `json:"高"`
	Label   string  `json:"标签"`
	Score   float64 `json:"精度"`
	CenterX int     `json:"-"` //中心坐标X
	CenterY int     `json:"-"` //中心坐标Y
}

// 模型初始化
func initialization() {
	labelPath := files.Path("./assets/label.txt")
	fileCheck(labelPath)

	dbParamPath := files.Path("./assets/db.param")
	fileCheck(dbParamPath)

	dbBinPath := files.Path("./assets/db.bin")
	fileCheck(dbBinPath)

	recParamPath := files.Path("./assets/rec.param")
	fileCheck(recParamPath)

	recBinPath := files.Path("./assets/rec.bin")
	fileCheck(recBinPath)

	pointer = C.newPpocr()
	if pointer == nil {
		_, _ = fmt.Fprintln(os.Stderr, "[ppocr] 创建实例失败")
		os.Exit(1)
	}

	cLabelPath := C.CString(labelPath)
	defer C.free(unsafe.Pointer(cLabelPath))

	cDbParamPath := C.CString(dbParamPath)
	defer C.free(unsafe.Pointer(cDbParamPath))

	cDbBinPath := C.CString(dbBinPath)
	defer C.free(unsafe.Pointer(cDbBinPath))

	cRecParamPath := C.CString(recParamPath)
	defer C.free(unsafe.Pointer(cRecParamPath))

	cRecBinPath := C.CString(recBinPath)
	defer C.free(unsafe.Pointer(cRecBinPath))

	result := C.loadModelPpocr(pointer, cLabelPath, cDbParamPath, cDbBinPath, cRecParamPath, cRecBinPath, 3)
	goResult := C.GoString(result)
	C.free(unsafe.Pointer(result))
	if goResult != "OK" {
		_, _ = fmt.Fprintln(os.Stderr, "[ppocr] "+goResult)
		os.Exit(1)
	}
}

// 在屏幕指定区域进行OCR文字识别
func Ocr(x1, y1, x2, y2 int, colorStr string) []Result {
	return ocr(images.CaptureScreen(x1, y1, x2, y2), x1, y1, colorStr)
}

// 从图像对象进行OCR文字识别
func OcrFromImage(img *image.NRGBA, colorStr string) []Result {
	return ocr(img, 0, 0, colorStr)
}

// 从Base64 编码识别
func OcrFromBase64(b64 string, colorStr string) []Result {
	return ocr(images.ReadFromBase64(b64), 0, 0, colorStr)
}

// 从文件路径识别
func OcrFromPath(path string, colorStr string) []Result {
	return ocr(images.ReadFromPath(path), 0, 0, colorStr)
}

// OCR识别函数
func ocr(img *image.NRGBA, x1, y1 int, colorStr string) []Result {
	mutex.Lock()
	defer mutex.Unlock()

	if pointer == nil {
		initialization()
	}

	if img == nil {
		return nil
	}

	colorStr = strings.ReplaceAll(colorStr, "-", ",")

	if colorStr != "" && !strings.Contains(colorStr, ",") {
		colorStr = colorStr + ",050505"
	}

	cColorStr := C.CString(colorStr)
	defer C.free(unsafe.Pointer(cColorStr))

	result := C.detectPpocr(
		pointer,
		(*C.uchar)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Dx()),
		C.int(img.Rect.Dy()),
		C.float(0.45),
		C.float(0.5),
		C.int(640),
		cColorStr,
	)
	goResult := C.GoString(result)
	C.free(unsafe.Pointer(result))

	var results []Result
	err := json.Unmarshal([]byte(goResult), &results)
	if err != nil {
		return nil
	}

	var validResults []Result
	for _, result := range results {
		if strings.TrimSpace(result.Label) == "" {
			continue
		}
		result.X += x1
		result.Y += y1
		result.CenterX = result.X + result.Width/2
		result.CenterY = result.Y + result.Height/2
		validResults = append(validResults, result)
	}

	return validResults
}

// 文件检查
func fileCheck(path string) {
	if !files.Exists(path) {
		_, _ = fmt.Fprintln(os.Stderr, "[ppocr] 缺少文件:"+path)
		os.Exit(1)
	}
}
