package yolo

// 目标检测库
/*
#include "yolo.h"
#include <stdlib.h>
#cgo arm64 LDFLAGS: -L../../resources/libs/arm64-v8a -lyolo
#cgo amd64 LDFLAGS: -L../../resources/libs/x86_64 -lyolo
#cgo 386 LDFLAGS: -L../../resources/libs/x86 -lyolo
*/
import "C" 
import (
	"encoding/json"
	"fmt"
	"image"
	"strings"
	"sync"
	"unsafe"

	"github.com/xiaocainiao633/Genie1.0--/images"
)

type Yolo struct {
	pointer *C.Yolo
	img     *image.NRGBA
	mutex   sync.Mutex
}

// 检测结果
type DetectResult struct {
	X       int     `json:"X"`
	Y       int     `json:"Y"`
	Width   int     `json:"宽"`
	Height  int     `json:"高"`
	Label   string  `json:"标签"`
	Score   float64 `json:"精度"`
	CenterX int     `json:"-"` //中心坐标X
	CenterY int     `json:"-"` //中心坐标Y
}

// new
func New(version string, cpuThreadNum int, paramPath, binPath, labels string) *Yolo {
	pointer := C.newYolo()
	if pointer == nil {
		return nil
	}

	cVersion := C.CString(strings.ToLower(version))
	defer C.free(unsafe.Pointer(cVersion))

	cParamPath := C.CString(paramPath)
	defer C.free(unsafe.Pointer(cParamPath))

	cBinPath := C.CString(binPath)
	defer C.free(unsafe.Pointer(cBinPath))

	cLabels := C.CString(labels)
	defer C.free(unsafe.Pointer(cLabels))

	result := C.loadModelYolo(pointer, cVersion, cParamPath, cBinPath, cLabels, C.int(cpuThreadNum))
	goResult := C.GoString(result)
	C.free(unsafe.Pointer(result))
	if goResult != "OK" {
		fmt.Println("yolo:" + goResult)
		return nil
	}
	return &Yolo{pointer: pointer}
}

// 设置待检测图像
func (y *Yolo) SetImage(img *image.NRGBA) {
	y.mutex.Lock()
	defer y.mutex.Unlock()
	y.img = img
}

// 执行目标检测
func (y *Yolo) Detect(x1, y1, x2, y2 int) []DetectResult {
	y.mutex.Lock()
	defer y.mutex.Unlock()
	var img *image.NRGBA
	if y.img != nil {
		img = y.getImage(x1, y1, x2, y2)
		defer func() {
			y.img = nil
		}()
	} else {
		img = images.CaptureScreen(x1, y1, x2, y2)
	}
	if img == nil {
		return nil
	}

	result := C.detectYolo(
		y.pointer,
		(*C.uchar)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Dx()),
		C.int(img.Rect.Dy()),
		C.float(0.5),
		C.float(0.45),
		C.int(640),
	)
	goResult := C.GoString(result)
	C.free(unsafe.Pointer(result))
	var results []DetectResult
	err := json.Unmarshal([]byte(goResult), &results)
	if err != nil {
		return nil
	}

	var validResults []DetectResult
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

// 释放资源
func (y *Yolo) Close() {
	C.closeYolo(y.pointer)
}

func (yl *Yolo) getImage(x1, y1, x2, y2 int) *image.NRGBA {
	width1 := yl.img.Rect.Dx()
	height1 := yl.img.Rect.Dy()

	if x2 == 0 {
		x2 = width1
	}
	if y2 == 0 {
		y2 = height1
	}

	if x1 < 0 || y1 < 0 || x2 > width1 || y2 > height1 || x1 >= x2 || y1 >= y2 {
		return nil
	}

	if x1 == 0 && y1 == 0 && x2 == width1 && y2 == height1 {
		return yl.img
	}
	width2 := x2 - x1
	height2 := y2 - y1
	img := image.NewNRGBA(image.Rect(0, 0, width2, height2))
	for y := 0; y < height2; y++ {
		srcOffset := (y1+y)*yl.img.Stride + x1*4
		destOffset := y * img.Stride
		copy(img.Pix[destOffset:], yl.img.Pix[srcOffset:srcOffset+width2*4])
	}
	return img
}
