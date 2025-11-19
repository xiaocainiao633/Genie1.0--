# AutoGo API 配合机制与自动化测试指南

## 目录

1. [框架概述](#框架概述)
2. [核心模块架构](#核心模块架构)
3. [API 配合机制](#api-配合机制)
4. [自动化测试脚本示例](#自动化测试脚本示例)
5. [最佳实践](#最佳实践)

---

## 框架概述

**AutoGo** 是一个高性能的 Android 自动化测试框架，基于 Go 语言开发。提供了完整的 Android 设备控制、图像识别、UI 操作等功能，支持跨应用自动化操作。

### 核心特性

- **无需安装 APK**：二进制文件可直接在 Android 系统上运行
- **跨应用操作**：可以操作任意 Android 应用
- **多种识别方式**：支持图像匹配、OCR、UI 控件识别、目标检测
- **完整的设备控制**：触摸、按键、应用管理、系统控制

---

## 核心模块架构

### 模块分类

AutoGo 的模块可以分为以下几个层次：

```
┌─────────────────────────────────────────────────┐
│           应用层 (Application Layer)             │
│  app, system, device, files, storages, https    │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│           交互层 (Interaction Layer)            │
│  motion, ime, uiacc                             │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│           识别层 (Recognition Layer)            │
│  opencv, ppocr, yolo, images                   │
└─────────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────────┐
│           基础层 (Foundation Layer)             │
│  utils, rhino, console, imgui                  │
└─────────────────────────────────────────────────┘
```

---

## API 配合机制

### 1. 图像识别 + 触摸操作

**典型流程**：截图 → 图像识别 → 获取坐标 → 点击

```go
// 模块配合：images + opencv + motion
import (
    "github.com/xiaocainiao633/Genie1.0--/images"
    "github.com/xiaocainiao633/Genie1.0--/opencv"
    "github.com/xiaocainiao633/Genie1.0--/motion"
    "os"
)

// 1. 读取模板图片
templateBytes, _ := os.ReadFile("./templates/button.png")

// 2. 在屏幕中查找图片
x, y := opencv.FindImage(0, 0, 0, 0, &templateBytes, false, 1.0, 0.8)

// 3. 如果找到，执行点击
if x != -1 && y != -1 {
    motion.Click(x, y, 1)
}
```

**配合说明**：
- `images.CaptureScreen()` 提供屏幕截图数据
- `opencv.FindImage()` 使用 OpenCV 进行模板匹配
- `motion.Click()` 执行实际的触摸操作

---

### 2. OCR 识别 + UI 操作

**典型流程**：截图 → OCR 识别文字 → 根据文字定位 → 操作

```go
// 模块配合：images + ppocr + motion
import (
    "github.com/xiaocainiao633/Genie1.0--/images"
    "github.com/xiaocainiao633/Genie1.0--/ppocr"
    "github.com/xiaocainiao633/Genie1.0--/motion"
    "strings"
)

// 1. 截取屏幕
img := images.CaptureScreen(0, 0, 0, 0)

// 2. OCR 识别文字
results := ppocr.OcrFromImage(img, "")

// 3. 查找特定文字并点击
for _, result := range results {
    if strings.Contains(result.Label, "登录") {
        motion.Click(result.CenterX, result.CenterY, 1)
        break
    }
}
```

**配合说明**：
- `images.CaptureScreen()` 提供图像数据
- `ppocr.OcrFromImage()` 识别图像中的文字
- `motion.Click()` 根据识别结果执行操作

---

### 3. UI 控件识别 + 自动化操作

**典型流程**：构建选择器 → 查找控件 → 执行操作

```go
// 模块配合：uiacc + motion
import (
    "github.com/xiaocainiao633/Genie1.0--/uiacc"
    "github.com/xiaocainiao633/Genie1.0--/motion"
)

// 1. 查找并点击文本为"确定"的按钮
obj := uiacc.New().Text("确定").FindOnce()
if obj != nil {
    obj.Click()
}

// 2. 等待元素出现
obj := uiacc.New().Id("login_button").WaitFor(5000)
if obj != nil {
    obj.Click()
}

// 3. 设置输入框文本
uiacc.New().Editable(true).FindOnce().SetText("Hello World")
```

**配合说明**：
- `uiacc.New()` 创建选择器
- `uiacc.FindOnce()` / `uiacc.WaitFor()` 查找控件
- `UiObject.Click()` / `UiObject.SetText()` 执行操作
- 内部使用 `motion` 模块执行实际触摸

---

### 4. 应用管理 + 自动化测试

**典型流程**：启动应用 → 执行测试 → 验证结果 → 清理

```go
// 模块配合：app + images + opencv + motion
import (
    "github.com/xiaocainiao633/Genie1.0--/app"
    "github.com/xiaocainiao633/Genie1.0--/opencv"
    "github.com/xiaocainiao633/Genie1.0--/motion"
    "os"
    "time"
)

// 1. 启动应用
app.Launch("com.example.app", 0)
time.Sleep(2000)

// 2. 验证应用是否启动成功
currentPkg := app.CurrentPackage()
if currentPkg != "com.example.app" {
    panic("应用启动失败")
}

// 3. 执行测试操作
templateBytes, _ := os.ReadFile("./templates/start_button.png")
x, y := opencv.FindImage(0, 0, 0, 0, &templateBytes, false, 1.0, 0.8)
if x != -1 {
    motion.Click(x, y, 1)
}

// 4. 清理：强制停止应用
app.ForceStop("com.example.app")
```

---

### 5. 多识别方式组合

**典型流程**：图像匹配失败 → OCR 识别 → UI 控件识别（降级策略）

```go
// 模块配合：opencv + ppocr + uiacc
import (
    "github.com/xiaocainiao633/Genie1.0--/opencv"
    "github.com/xiaocainiao633/Genie1.0--/ppocr"
    "github.com/xiaocainiao633/Genie1.0--/uiacc"
    "github.com/xiaocainiao633/Genie1.0--/images"
    "os"
    "strings"
)

func FindAndClickButton(buttonText string) bool {
    // 策略1：图像模板匹配
    templateBytes, err := os.ReadFile("./templates/" + buttonText + ".png")
    if err == nil {
        x, y := opencv.FindImage(0, 0, 0, 0, &templateBytes, false, 1.0, 0.8)
        if x != -1 && y != -1 {
            motion.Click(x, y, 1)
            return true
        }
    }
    
    // 策略2：OCR 文字识别
    img := images.CaptureScreen(0, 0, 0, 0)
    results := ppocr.OcrFromImage(img, "")
    for _, result := range results {
        if strings.Contains(result.Label, buttonText) {
            motion.Click(result.CenterX, result.CenterY, 1)
            return true
        }
    }
    
    // 策略3：UI 控件识别
    obj := uiacc.New().Text(buttonText).FindOnce()
    if obj != nil {
        obj.Click()
        return true
    }
    
    return false
}
```

---

### 6. 数据持久化 + 测试配置

**典型流程**：保存测试配置 → 读取配置 → 执行测试 → 保存结果

```go
// 模块配合：storages + files
import (
    "github.com/xiaocainiao633/Genie1.0--/storages"
    "github.com/xiaocainiao633/Genie1.0--/files"
    "encoding/json"
)

// 1. 保存测试配置
config := map[string]string{
    "app_package": "com.example.app",
    "test_timeout": "30",
}
configJson, _ := json.Marshal(config)
storages.Put("test_config", "app_config", string(configJson))

// 2. 读取配置
configStr := storages.Get("test_config", "app_config")
json.Unmarshal([]byte(configStr), &config)

// 3. 保存测试结果
result := map[string]interface{}{
    "test_name": "登录测试",
    "passed": true,
    "duration": 5.2,
}
resultJson, _ := json.Marshal(result)
files.Write("./test_results.json", string(resultJson))
```

---

### 7. 系统监控 + 性能测试

**典型流程**：启动应用 → 监控资源 → 执行操作 → 分析性能

```go
// 模块配合：app + system + device
import (
    "github.com/xiaocainiao633/Genie1.0--/app"
    "github.com/xiaocainiao633/Genie1.0--/system"
    "github.com/xiaocainiao633/Genie1.0--/device"
    "time"
)

// 1. 获取设备信息
fmt.Printf("设备型号: %s\n", device.Model)
fmt.Printf("Android版本: %s\n", device.Release)
fmt.Printf("屏幕分辨率: %dx%d\n", device.Width, device.Height)

// 2. 启动应用
app.Launch("com.example.app", 0)
time.Sleep(2000)

// 3. 获取应用 PID
pid := system.GetPid("com.example.app")

// 4. 监控资源使用
memory := system.GetMemoryUsage(pid)
cpu := system.GetCpuUsage(pid)
fmt.Printf("内存使用: %d KB\n", memory)
fmt.Printf("CPU使用率: %.2f%%\n", cpu)
```

---

## 自动化测试脚本示例

### 示例 1：登录界面自动化测试

```go
package main

import (
    "fmt"
    "os"
    "strings"
    "time"
    
    "github.com/xiaocainiao633/Genie1.0--/app"
    "github.com/xiaocainiao633/Genie1.0--/images"
    "github.com/xiaocainiao633/Genie1.0--/motion"
    "github.com/xiaocainiao633/Genie1.0--/opencv"
    "github.com/xiaocainiao633/Genie1.0--/ppocr"
    "github.com/xiaocainiao633/Genie1.0--/uiacc"
    "github.com/xiaocainiao633/Genie1.0--/utils"
)

// TestResult 测试结果结构
type TestResult struct {
    TestName string
    Passed   bool
    Message  string
    Duration time.Duration
}

// TestLoginInterface 测试登录界面
func TestLoginInterface() TestResult {
    startTime := time.Now()
    testName := "登录界面测试"
    
    // 步骤1：启动应用
    fmt.Println("[步骤1] 启动应用...")
    if !app.Launch("com.example.app", 0) {
        return TestResult{
            TestName: testName,
            Passed:   false,
            Message:  "应用启动失败",
            Duration: time.Since(startTime),
        }
    }
    utils.Sleep(3000)
    
    // 步骤2：验证是否进入登录页面
    fmt.Println("[步骤2] 验证登录页面...")
    loginPageTemplate, err := os.ReadFile("./templates/login_page.png")
    if err == nil {
        x, y := opencv.FindImage(0, 0, 0, 0, &loginPageTemplate, false, 1.0, 0.8)
        if x == -1 {
            return TestResult{
                TestName: testName,
                Passed:   false,
                Message:  "未检测到登录页面",
                Duration: time.Since(startTime),
            }
        }
    }
    
    // 步骤3：查找用户名输入框
    fmt.Println("[步骤3] 查找用户名输入框...")
    usernameInput := uiacc.New().Editable(true).Index(0).FindOnce()
    if usernameInput == nil {
        // 降级策略：使用 OCR 查找
        img := images.CaptureScreen(0, 0, 0, 0)
        results := ppocr.OcrFromImage(img, "")
        found := false
        for _, result := range results {
            if strings.Contains(result.Label, "用户名") || strings.Contains(result.Label, "账号") {
                motion.Click(result.CenterX, result.CenterY-50, 1) // 点击输入框上方
                found = true
                break
            }
        }
        if !found {
            return TestResult{
                TestName: testName,
                Passed:   false,
                Message:  "未找到用户名输入框",
                Duration: time.Since(startTime),
            }
        }
    } else {
        usernameInput.Click()
    }
    utils.Sleep(500)
    
    // 步骤4：输入用户名
    fmt.Println("[步骤4] 输入用户名...")
    usernameInput.SetText("testuser")
    utils.Sleep(1000)
    
    // 步骤5：查找密码输入框
    fmt.Println("[步骤5] 查找密码输入框...")
    passwordInput := uiacc.New().Editable(true).Index(1).FindOnce()
    if passwordInput == nil {
        return TestResult{
            TestName: testName,
            Passed:   false,
            Message:  "未找到密码输入框",
            Duration: time.Since(startTime),
        }
    }
    passwordInput.Click()
    utils.Sleep(500)
    
    // 步骤6：输入密码
    fmt.Println("[步骤6] 输入密码...")
    passwordInput.SetText("password123")
    utils.Sleep(1000)
    
    // 步骤7：点击登录按钮
    fmt.Println("[步骤7] 点击登录按钮...")
    loginButton := uiacc.New().Text("登录").FindOnce()
    if loginButton == nil {
        // 降级策略：图像匹配
        loginBtnTemplate, err := os.ReadFile("./templates/login_button.png")
        if err == nil {
            x, y := opencv.FindImage(0, 0, 0, 0, &loginBtnTemplate, false, 1.0, 0.8)
            if x != -1 {
                motion.Click(x, y, 1)
            } else {
                return TestResult{
                    TestName: testName,
                    Passed:   false,
                    Message:  "未找到登录按钮",
                    Duration: time.Since(startTime),
                }
            }
        } else {
            return TestResult{
                TestName: testName,
                Passed:   false,
                Message:  "未找到登录按钮",
                Duration: time.Since(startTime),
            }
        }
    } else {
        loginButton.Click()
    }
    utils.Sleep(3000)
    
    // 步骤8：验证登录结果
    fmt.Println("[步骤8] 验证登录结果...")
    
    // 检测错误弹窗
    errorTemplates := []string{
        "./templates/error_dialog.png",
        "./templates/login_failed.png",
    }
    for _, templatePath := range errorTemplates {
        templateBytes, err := os.ReadFile(templatePath)
        if err == nil {
            x, y := opencv.FindImage(0, 0, 0, 0, &templateBytes, false, 1.0, 0.7)
            if x != -1 {
                return TestResult{
                    TestName: testName,
                    Passed:   false,
                    Message:  "登录失败：检测到错误弹窗",
                    Duration: time.Since(startTime),
                }
            }
        }
    }
    
    // 检测成功标志
    img := images.CaptureScreen(0, 0, 0, 0)
    results := ppocr.OcrFromImage(img, "")
    successKeywords := []string{"欢迎", "首页", "主页", "成功"}
    for _, result := range results {
        for _, keyword := range successKeywords {
            if strings.Contains(result.Label, keyword) {
                return TestResult{
                    TestName: testName,
                    Passed:   true,
                    Message:  "登录成功",
                    Duration: time.Since(startTime),
                }
            }
        }
    }
    
    // 使用 UI 控件验证
    homeElement := uiacc.New().PackageName("com.example.app").FindOnce()
    if homeElement != nil {
        return TestResult{
            TestName: testName,
            Passed:   true,
            Message:  "登录成功（通过UI验证）",
            Duration: time.Since(startTime),
        }
    }
    
    return TestResult{
        TestName: testName,
        Passed:   false,
        Message:  "无法确定登录状态",
        Duration: time.Since(startTime),
    }
}

// TestErrorHandling 测试错误处理
func TestErrorHandling() TestResult {
    startTime := time.Now()
    testName := "错误处理测试"
    
    // 测试：输入错误密码
    fmt.Println("[测试] 输入错误密码...")
    
    // 输入错误的用户名和密码
    usernameInput := uiacc.New().Editable(true).Index(0).FindOnce()
    if usernameInput != nil {
        usernameInput.SetText("wronguser")
    }
    utils.Sleep(500)
    
    passwordInput := uiacc.New().Editable(true).Index(1).FindOnce()
    if passwordInput != nil {
        passwordInput.SetText("wrongpass")
    }
    utils.Sleep(500)
    
    // 点击登录
    loginButton := uiacc.New().Text("登录").FindOnce()
    if loginButton != nil {
        loginButton.Click()
    }
    utils.Sleep(2000)
    
    // 验证是否出现错误提示
    img := images.CaptureScreen(0, 0, 0, 0)
    results := ppocr.OcrFromImage(img, "")
    errorKeywords := []string{"错误", "失败", "不正确", "无效"}
    for _, result := range results {
        for _, keyword := range errorKeywords {
            if strings.Contains(result.Label, keyword) {
                return TestResult{
                    TestName: testName,
                    Passed:   true,
                    Message:  "错误处理正常",
                    Duration: time.Since(startTime),
                }
            }
        }
    }
    
    return TestResult{
        TestName: testName,
        Passed:   false,
        Message:  "未检测到错误提示",
        Duration: time.Since(startTime),
    }
}

// TestUIElements 测试UI元素存在性
func TestUIElements() TestResult {
    startTime := time.Now()
    testName := "UI元素存在性测试"
    
    requiredElements := []struct {
        name     string
        selector *uiacc.Uiacc
    }{
        {"用户名输入框", uiacc.New().Editable(true).Index(0)},
        {"密码输入框", uiacc.New().Editable(true).Index(1)},
        {"登录按钮", uiacc.New().Text("登录")},
    }
    
    missingElements := []string{}
    for _, element := range requiredElements {
        obj := element.selector.FindOnce()
        if obj == nil {
            missingElements = append(missingElements, element.name)
        }
    }
    
    if len(missingElements) > 0 {
        return TestResult{
            TestName: testName,
            Passed:   false,
            Message:  fmt.Sprintf("缺失元素: %v", missingElements),
            Duration: time.Since(startTime),
        }
    }
    
    return TestResult{
        TestName: testName,
        Passed:   true,
        Message:  "所有UI元素存在",
        Duration: time.Since(startTime),
    }
}

// RunAllTests 运行所有测试
func RunAllTests() []TestResult {
    results := []TestResult{}
    
    // 运行测试用例
    results = append(results, TestUIElements())
    results = append(results, TestErrorHandling())
    results = append(results, TestLoginInterface())
    
    return results
}

// GenerateReport 生成测试报告
func GenerateReport(results []TestResult) {
    fmt.Println("\n" + strings.Repeat("=", 60))
    fmt.Println("自动化测试报告")
    fmt.Println(strings.Repeat("=", 60))
    
    passedCount := 0
    failedCount := 0
    totalDuration := time.Duration(0)
    
    for _, result := range results {
        status := "❌ 失败"
        if result.Passed {
            status = "✅ 通过"
            passedCount++
        } else {
            failedCount++
        }
        
        fmt.Printf("\n[%s] %s\n", status, result.TestName)
        fmt.Printf("  消息: %s\n", result.Message)
        fmt.Printf("  耗时: %v\n", result.Duration)
        totalDuration += result.Duration
    }
    
    fmt.Println("\n" + strings.Repeat("-", 60))
    fmt.Printf("总计: %d 个测试\n", len(results))
    fmt.Printf("通过: %d 个\n", passedCount)
    fmt.Printf("失败: %d 个\n", failedCount)
    fmt.Printf("总耗时: %v\n", totalDuration)
    fmt.Println(strings.Repeat("=", 60))
}

func main() {
    fmt.Println("开始执行自动化测试...")
    
    // 运行所有测试
    results := RunAllTests()
    
    // 生成报告
    GenerateReport(results)
    
    // 清理：强制停止应用
    app.ForceStop("com.example.app")
}
```

---

### 示例 2：购物应用自动化测试

```go
package main

import (
    "fmt"
    "os"
    "strings"
    "time"
    
    "github.com/xiaocainiao633/Genie1.0--/app"
    "github.com/xiaocainiao633/Genie1.0--/images"
    "github.com/xiaocainiao633/Genie1.0--/motion"
    "github.com/xiaocainiao633/Genie1.0--/opencv"
    "github.com/xiaocainiao633/Genie1.0--/ppocr"
    "github.com/xiaocainiao633/Genie1.0--/uiacc"
    "github.com/xiaocainiao633/Genie1.0--/utils"
)

// TestShoppingFlow 测试购物流程
func TestShoppingFlow() {
    // 1. 启动购物应用
    app.Launch("com.shopping.app", 0)
    utils.Sleep(3000)
    
    // 2. 搜索商品（使用 OCR 识别搜索框）
    img := images.CaptureScreen(0, 0, 0, 0)
    results := ppocr.OcrFromImage(img, "")
    for _, result := range results {
        if strings.Contains(result.Label, "搜索") {
            motion.Click(result.CenterX, result.CenterY, 1)
            break
        }
    }
    utils.Sleep(1000)
    
    // 3. 输入搜索关键词
    searchInput := uiacc.New().Editable(true).Focused(true).FindOnce()
    if searchInput != nil {
        searchInput.SetText("手机")
    }
    utils.Sleep(1000)
    
    // 4. 点击搜索按钮（图像匹配）
    searchBtnTemplate, _ := os.ReadFile("./templates/search_button.png")
    x, y := opencv.FindImage(0, 0, 0, 0, &searchBtnTemplate, false, 1.0, 0.8)
    if x != -1 {
        motion.Click(x, y, 1)
    }
    utils.Sleep(3000)
    
    // 5. 选择第一个商品（UI控件识别）
    firstProduct := uiacc.New().Clickable(true).Index(0).FindOnce()
    if firstProduct != nil {
        firstProduct.Click()
    }
    utils.Sleep(2000)
    
    // 6. 添加到购物车（多种方式组合）
    addToCartBtn := uiacc.New().Text("加入购物车").FindOnce()
    if addToCartBtn == nil {
        // 降级：图像匹配
        cartBtnTemplate, _ := os.ReadFile("./templates/add_to_cart.png")
        x, y := opencv.FindImage(0, 0, 0, 0, &cartBtnTemplate, false, 1.0, 0.8)
        if x != -1 {
            motion.Click(x, y, 1)
        }
    } else {
        addToCartBtn.Click()
    }
    utils.Sleep(2000)
    
    // 7. 验证是否添加成功（OCR 识别提示信息）
    img = images.CaptureScreen(0, 0, 0, 0)
    results = ppocr.OcrFromImage(img, "")
    for _, result := range results {
        if strings.Contains(result.Label, "成功") || strings.Contains(result.Label, "已添加") {
            fmt.Println("✅ 商品已成功添加到购物车")
            return
        }
    }
    
    fmt.Println("❌ 无法确认商品是否添加成功")
}
```

---

## 最佳实践

### 1. 识别方式选择策略

```
优先级1: UI控件识别 (uiacc)
  ↓ 失败
优先级2: 图像模板匹配 (opencv)
  ↓ 失败
优先级3: OCR文字识别 (ppocr)
  ↓ 失败
优先级4: 目标检测 (yolo)
```

### 2. 错误处理

- **超时处理**：使用 `WaitFor()` 等待元素出现
- **降级策略**：多种识别方式组合使用
- **异常检测**：定期检查错误弹窗
- **重试机制**：失败后自动重试

### 3. 性能优化

- **缓存截图**：避免频繁截图
- **区域限制**：只在必要区域进行识别
- **灰度匹配**：使用 `isGray=true` 提升速度
- **相似度调整**：根据场景调整相似度阈值

### 4. 测试数据管理

- **使用 storages**：保存测试配置和结果
- **使用 files**：保存测试日志和截图
- **结构化数据**：使用 JSON 格式存储

### 5. 代码组织

```go
// 推荐的项目结构
project/
├── main.go              // 主程序
├── tests/               // 测试用例
│   ├── login_test.go
│   └── shopping_test.go
├── templates/           // 图像模板
│   ├── login_button.png
│   └── search_button.png
├── config/              // 配置文件
│   └── app_config.json
└── results/             // 测试结果
    └── test_report.json
```

---

## 总结

AutoGo 框架通过模块化的设计，提供了灵活的 API 组合方式。通过合理组合不同的识别和操作模块，可以构建强大的自动化测试脚本。关键是要理解各个模块的职责和配合方式，根据实际场景选择最合适的组合策略。

---
