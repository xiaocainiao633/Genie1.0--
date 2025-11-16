package ime

// 用于控制输入法行为，实现文本输入等功能
// 工作流程：
// 1. 获取剪切板文本
// 2. 设置剪切板文本
// 3. 输入文本
// 4. 模拟按键
// 5. 获取输入法列表
// 6. 设置当前输入法

// 获取剪切板文本
// text := ime.GetClipText()

// 设置剪切板文本
// ime.SetClipText("Hello, World!")

// 输入文本
// ime.InputText("Hello, World!")

import (
	"bufio"
	_ "embed" // 嵌入JavaScript代码
	"encoding/base64"
	"github.com/xiaocainiao633/Genie1.0--/motion"
	"github.com/xiaocainiao633/Genie1.0--/rhino"
	"github.com/xiaocainiao633/Genie1.0--/utils"
	"strconv"
	"strings"
	"time"
)

//go:embed ime.js
var _ime string

func init() {
	rhino.Eval("_ime", _ime) // 初始化JavaScript引擎
}

// GetClipText 获取剪切板内容
func GetClipText() string {
	return rhino.Eval("_ime", "getClipboardText()")
}

// SetClipText 设置剪切板内容
func SetClipText(text string) bool {
	if text != "" {
		text = base64.StdEncoding.EncodeToString([]byte(text))
	}
	return rhino.Eval("_ime", "setClipboardText('"+text+"')") == "true"
}

// InputText 输入文本
func InputText(text string) {
	// 判断字符串是否包含中文
	if containsChinese(text) {
		SetClipText(text)
		motion.KeyAction(279)
		time.Sleep(time.Duration(50) * time.Millisecond)
		SetClipText("")
	} else {
		// 如果不包含中文，则使用 Shell 输入文本
		utils.Shell("input text " + strconv.Quote(text))
	}
}

// KeyAction 模拟按键
func KeyAction(code int) {
	utils.Shell("am broadcast -a com.autogo --es message keyaction --es data " + strconv.Itoa(code))
}

// GetIMEList 获取输入法列表
func GetIMEList() []string {
	shellOutput := utils.Shell("ime list -a | grep mId")

	var imeList []string

	scanner := bufio.NewScanner(strings.NewReader(shellOutput))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "mId=") {
			mId := strings.SplitN(line, " ", 2)[0]
			mId = strings.TrimPrefix(mId, "mId=")
			imeList = append(imeList, mId)
		}
	}

	if err := scanner.Err(); err != nil {
		//fmt.Println("Scanner error:", err)
		return nil
	}

	return imeList
}

// SetCurrentIME 设置当前输入法
func SetCurrentIME(imeID string) {
	utils.Shell("ime enable " + imeID)
	utils.Shell("ime set " + imeID)
}

// containsChinese 检查字符串是否包含中文
func containsChinese(s string) bool {
	for _, r := range s {
		if r >= '\u4e00' && r <= '\u9fff' { // 中文字符范围
			return true
		}
	}
	return false
}
