package utils

/*
#include <sys/mman.h>
#include <stdint.h>
#include <stdlib.h>
#cgo arm64 LDFLAGS: -L../libs/arm64-v8a -lAutoGo
#cgo amd64 LDFLAGS: -L../libs/x86_64 -lAutoGo
#cgo 386 LDFLAGS: -L../libs/x86 -lAutoGo

int shmem_init(int size);
int shmem_write(char* buffer, size_t data_size);
int shmem_read(void *buffer, int buf_size);
*/
import "C"
import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

//go:embed utils.dex
var _utils []byte

//go:embed libashmem_arm64.so
var _libashmem_arm64 []byte

//go:embed libashmem_amd.so
var _libashmem_amd []byte

var conn net.Conn
var msgChanMap = make(map[string]chan string)
var mu sync.Mutex
var currentMsgId int
var ashmemSize int

//go:embed utils.js
var _utils_js string

func init() {
	var err error
	addr := net.UnixAddr{
		Name: "@ags.socket",
		Net:  "unix",
	}
	conn, err = net.DialUnix("unix", nil, &addr)
	if err != nil {
		path := filepath.Dir(os.Args[0])
		Shell("pkill -f com.ags.Main")
		var appProcess, library string
		if runtime.GOARCH == "arm64" {
			_ = os.WriteFile(path+"/libashmem.so", _libashmem_arm64, 0644)
			appProcess = "app_process"
			library = "export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/system/lib64:/system/lib;"
		} else {
			_ = os.WriteFile(path+"/libashmem.so", _libashmem_amd, 0644)
			appProcess = "app_process32"
			library = "export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/system/lib;"
		}
		_ = os.WriteFile(path+"/ags.dex", _utils, 0644)

		command := "nohup sh -c '" + library + appProcess + " -Djava.class.path=" + filepath.Dir(os.Args[0]) + "/ags.dex / com.ags.Main ' > /dev/null 2>&1 &"

		cCommand := C.CString(command)
		defer C.free(unsafe.Pointer(cCommand))

		C.system(cCommand)

		for i := 0; i < 151; i++ {
			Sleep(50)
			conn, err = net.DialUnix("unix", nil, &addr)
			if err == nil {
				break
			}
			if i == 150 {
				handleError(err)
			}
		}
	}
	go func() {
		reader := bufio.NewReader(conn)
		for {
			// 先读取消息ID和消息长度，共12个字节
			header := make([]byte, 12)
			_, err = reader.Read(header)
			if err != nil {
				fmt.Println("Error reading message header:", err)
				return
			}
			// 解析消息ID（前6个字节）
			msgId := string(header[:6])

			// 解析消息长度（后6个字节）
			lengthStr := string(header[6:12])
			msgLength, err := strconv.Atoi(lengthStr)
			if err != nil {
				fmt.Println("Error parsing message length:", err)
				return
			}

			// 读取消息体
			message := make([]byte, msgLength)
			_, err = io.ReadFull(reader, message)
			if err != nil {
				fmt.Println("Error reading message body:", err)
				return
			}

			// 将消息发送到对应的通道
			mu.Lock()
			if ch, ok := msgChanMap[msgId]; ok {
				ch <- string(message)
				close(ch)
				delete(msgChanMap, msgId)
			}
			mu.Unlock()
		}
	}()
	parts := strings.Split(CallJavaMethod("js", "_utils|"+_utils_js), "x")
	var w, h int
	if len(parts) == 2 {
		w, _ = strconv.Atoi(parts[0])
		h, _ = strconv.Atoi(parts[1])
	} else {
		handleError(fmt.Errorf("设备分辨率获取失败"))
	}
	ashmemSize = (w + 40) * (h + 40) * 4
	if int(C.shmem_init(C.int(ashmemSize))) < 0 {
		handleError(fmt.Errorf("共享内存映射失败"))
	}
}

func CallJavaMethod(model string, str string) string {
	respChan := make(chan string)
	mu.Lock()
	currentMsgId++
	if currentMsgId > 999999 {
		currentMsgId = 1
	}
	msgId := fmt.Sprintf("%06d", currentMsgId)
	_, _ = conn.Write([]byte(model + "|" + msgId + "|" + str + "\u001E"))
	msgChanMap[msgId] = respChan
	mu.Unlock()
	select {
	case resp := <-respChan:
		mu.Lock()
		delete(msgChanMap, msgId) // 收到响应后清理 map
		mu.Unlock()
		return resp
	case <-time.After(12 * time.Second):
		mu.Lock()
		delete(msgChanMap, msgId) // 超时后清理 map
		mu.Unlock()
		fmt.Println("Timeout waiting for response " + str)
		return ""
	}
}

func Send(str string) {
	mu.Lock()
	defer mu.Unlock()
	_, _ = conn.Write([]byte(str + "\u001E"))
}

func GetBitMapData() ([]byte, error) {
	return shmemRead()
}

func shmemRead() ([]byte, error) {
	// 第一次读取时，传入 nil 来获取数据长度
	dataLen := int(C.shmem_read(nil, 0)) // 先获取数据长度
	if dataLen < 0 {
		return nil, fmt.Errorf("共享内存读取失败")
	}

	if dataLen == 0 {
		return nil, nil
	}

	// 创建一个大小为 dataLen 的缓冲区
	buffer := make([]byte, dataLen)

	// 读取数据到缓冲区，传入 clearData 参数
	bufferPtr := unsafe.Pointer(&buffer[0])
	readSize := C.shmem_read(bufferPtr, C.int(len(buffer)))
	if int(readSize) != dataLen {
		return nil, fmt.Errorf("共享内存数据读取错误 %d, got %d", dataLen, int(readSize))
	}

	return buffer, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[AutoGo] 出现错误: %v\n", err)
		os.Exit(1)
	}
}

func WmSize() (int, int) {
	parts := strings.Split(CallJavaMethod("js", "_utils|wmSize()"), "x")
	if len(parts) == 2 {
		w, _ := strconv.Atoi(parts[0])
		h, _ := strconv.Atoi(parts[1])
		return w, h
	}
	return 0, 0
}
