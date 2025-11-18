# RAG 对话系统使用指南

## 快速开始

### 步骤 1: 初始化知识库

```bash
cd AutoGo
go run agent/main.go -init -kb ./knowledge_base.db
```

这将创建SQLite知识库并填充默认的API文档。

### 步骤 2: 启动对话系统

```bash
go run main.go -kb ./knowledge_base.db -workspace ./workspace
```

### 步骤 3: 开始对话

```
你: 点击登录按钮

[系统] 正在处理您的请求...
[Agent] 正在分析用户需求...
[Agent] 代码生成完成
生成的代码:
---
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xiaocainiao633/Genie1.0--/uiacc"
	"github.com/xiaocainiao633/Genie1.0--/motion"
	"github.com/xiaocainiao633/Genie1.0--/utils"
)

func main() {
	fmt.Println("开始执行自动化测试...")

	// 查找并点击按钮
	obj := uiacc.New().Text("登录").FindOnce()
	if obj != nil {
		obj.Click()
		fmt.Println("点击成功")
	} else {
		fmt.Println("未找到按钮")
	}

	fmt.Println("测试完成")
}
---
```

## 使用示例

### 示例 1: 简单点击

**输入**: `点击登录按钮`

**生成代码**: 使用 `uiacc.New().Text("登录").FindOnce().Click()`

### 示例 2: 坐标点击

**输入**: `在坐标(100, 200)点击`

**生成代码**: `motion.Click(100, 200, 1)`

### 示例 3: 输入文本

**输入**: `在输入框中输入'Hello World'`

**生成代码**:
```go
inputObj := uiacc.New().Editable(true).FindOnce()
if inputObj != nil {
    inputObj.SetText("Hello World")
}
```

### 示例 4: 图像匹配点击

**输入**: `点击图片button.png`

**生成代码**:
```go
templateBytes, _ := os.ReadFile("button.png")
x, y := opencv.FindImage(0, 0, 0, 0, &templateBytes, false, 1.0, 0.8)
if x != -1 && y != -1 {
    motion.Click(x, y, 1)
}
```

### 示例 5: 文本验证

**输入**: `验证页面是否存在'主页'文字`

**生成代码**:
```go
img := images.CaptureScreen(0, 0, 0, 0)
results := ppocr.OcrFromImage(img, "")
for _, result := range results {
    if strings.Contains(result.Label, "主页") {
        fmt.Println("✅ 断言通过：文本存在")
        break
    }
}
```

### 示例 6: 图像验证

**输入**: `验证是否存在图片logo.png`

**生成代码**:
```go
templateBytes, _ := os.ReadFile("logo.png")
x, y := opencv.FindImage(0, 0, 0, 0, &templateBytes, false, 1.0, 0.8)
if x != -1 && y != -1 {
    fmt.Println("✅ 断言通过：图像存在")
} else {
    fmt.Println("❌ 断言失败：图像不存在")
}
```

### 示例 7: 启动应用

**输入**: `启动应用com.example.app`

**生成代码**:
```go
if app.Launch("com.example.app", 0) {
    fmt.Println("应用启动成功")
    utils.Sleep(2000)
}
```

### 示例 8: 等待操作

**输入**: `等待5秒`

**生成代码**: `utils.Sleep(5000)`

**输入**: `等待登录按钮出现`

**生成代码**:
```go
obj := uiacc.New().Text("登录").WaitFor(5000)
if obj != nil {
    fmt.Println("元素已出现")
}
```

### 示例 9: 滑动操作

**输入**: `从(100,200)滑动到(300,400)`

**生成代码**: `motion.Swipe(100, 200, 300, 400, 500)`

## 复杂场景示例

### 场景 1: 登录流程

**输入**: `测试登录流程：启动应用com.example.app，输入用户名'testuser'，输入密码'password123'，点击登录按钮，验证是否出现'欢迎'文字`

系统会生成完整的登录测试代码。

### 场景 2: 购物流程

**输入**: `搜索商品'手机'，选择第一个商品，添加到购物车`

系统会生成购物流程测试代码。

## 单次查询模式

不需要交互式对话时，可以使用单次查询：

```bash
go run agent/main.go \
  -kb ./knowledge_base.db \
  -workspace ./workspace \
  -query "点击登录按钮"
```

## API 查询

查询可用的API：

```go
agent := NewAgent(...)
apis, err := agent.GetAPIInfo("点击")
// 返回相关的API文档
```

## 自定义知识库

### 添加新的API文档

```go
doc := agent.APIDoc{
    Module:      "custom",
    Function:    "CustomFunction",
    Description: "自定义功能",
    Signature:   "func CustomFunction()",
    Parameters:  "无",
    Return:      "无",
    Example:     "custom.CustomFunction()",
    Keywords:    "自定义 custom",
}

kb.AddAPI(doc)
```

## 故障排除

### 问题 1: 知识库初始化失败

**解决方案**: 确保有写入权限，删除旧的知识库文件重新初始化。

### 问题 2: 代码编译失败

**检查**:
- Go环境是否正确
- 模块路径是否正确
- 依赖是否安装

### 问题 3: 意图识别不准确

**解决方案**: 
- 使用更明确的描述
- 参考示例修改查询
- 可以手动指定模块和操作类型

## 最佳实践

1. **明确描述**: 使用清晰、具体的描述
2. **指定模块**: 如果知道使用的模块，在查询中提及
3. **提供参数**: 对于需要参数的操作，明确提供参数值
4. **分步操作**: 复杂流程可以分步描述

## 扩展开发

参考 `ARCHITECTURE.md` 了解如何扩展系统功能。

## 自动化测试如何运行？

1. **自然语言 + 多轮记忆**：对话系统接收用户需求，并拼接最近的上下文（如上一条测试结果）
2. **RAG 检索**：系统到 SQLite + 向量知识库中检索符合需求的 API 文档 / 示例脚本
3. **脚本生成**：
   - 常规操作走规则模板（例如点击、输入、等待等）
   - 复杂业务（多步校验、流程控制）会把检索结果 + 历史上下文发送给 Ollama `llama3.2` 生成完整 Go 代码
4. **脚本编译**：把代码写入 `workspace` 并 `go build` 出可执行文件
5. **自动执行（可选）**：若启用 `-auto-exec`，系统会用 ADB 将二进制推送到 Android 设备执行
6. **测试报告**：收集编译日志、运行日志、错误堆栈，生成 JSON & Markdown 报告，包含“每一步的输入/输出/异常”

> 因此，项目是通过“脚本 + API”的方式进行自动化测试：AutoGo 会自适应调用 `motion`, `uiacc`, `opencv`, `ppocr`, `yolo` 等 API 来完成操作与断言。

## 如何判断页面是否有异常？

AutoGo 会组合多种视觉与 UI 技术对页面进行体检：

| 技术 | 调用 API | 作用 | 可检测异常 |
|------|----------|------|------------|
| 模板匹配 | `opencv.FindImage` | 比对界面截图与模板 | 按钮缺失、布局错位、主题颜色异常 |
| 目标检测 | `yolo.Detect` | 识别动态元素 | 弹窗、广告、遮挡层、异常浮窗 |
| OCR 识别 | `ppocr.Ocr` / `OcrFromImage` | 读取文本提示 | “错误”“失败”“网络异常”等字样 |
| 颜色检测 | `images.FindColor` | 校验状态条颜色 | 加载/成功/失败状态条颜色不符合预期 |
| UI 控件 | `uiacc.New().Text(...).FindOnce()` | 精准定位控件状态 | 控件不可见/不可点击/值错误 |

系统通常以“模板匹配 → YOLO → OCR → UI 属性”多级校验方式判定页面是否正常，并将判定过程写入测试报告（包括截图路径和异常说明）。

