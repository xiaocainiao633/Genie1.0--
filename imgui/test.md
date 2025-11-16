## 界面元素存在性测试

```
func TestLoginButtonExists() {
    gui.Init(false)
    
    // 1. 截图并OCR识别
    screenshot := takeScreenshot()
    textResults := ocrEngine.Recognize(screenshot)
    
    // 2. 查找"登录"按钮
    loginButtonFound := false
    for _, result := range textResults {
        if result.Text == "登录" {
            // 3. 可视化标记找到的元素
            gui.DrawRect(result.X1, result.Y1, result.X2, result.Y2, "#00FF00")
            gui.Toast("找到登录按钮")
            loginButtonFound = true
            break
        }
    }
    
    // 4. 测试断言
    if !loginButtonFound {
        gui.Alert("测试失败", "未找到登录按钮")
    }
}
```

## 完整业务流程测试

```
func TestUserLoginFlow() {
    gui.Init(false)
    gui.HudInit(10, 10, 400, 60, "#2B2B2B", 16)
    
    // 步骤1：点击登录入口
    gui.Toast("步骤1: 点击登录入口")
    tap(520, 880)  // 模拟点击登录按钮坐标
    utils.Sleep(2000)
    
    // 步骤2：输入用户名
    gui.Toast("步骤2: 输入用户名")
    tap(300, 600)  // 点击用户名输入框
    inputText("testuser")
    gui.HudSetText([]gui.TextItem{
        {TextColor: "#FFFF00", Text: "状态: 已输入用户名"},
    })
    
    // 步骤3：输入密码
    gui.Toast("步骤3: 输入密码")  
    tap(300, 700)  // 点击密码输入框
    inputText("password123")
    
    // 步骤4：点击登录并验证结果
    gui.Toast("步骤4: 执行登录")
    tap(400, 800)  // 点击登录按钮
    utils.Sleep(3000)
    
    // 验证登录成功
    screenshot := takeScreenshot()
    if ocrContains(screenshot, "欢迎") {
        gui.Toast("✅ 登录测试通过")
        gui.DrawRect(100, 100, 500, 200, "#00FF00")  // 绿色成功标记
    } else {
        gui.Toast("❌ 登录测试失败")
        gui.DrawRect(100, 100, 500, 200, "#FF0000")  // 红色失败标记
    }
}
```

## 数据验证场景测试

```
func TestBalanceDisplay() {
    gui.Init(false)
    
    // 1. 导航到余额页面
    tap(150, 150)  // 点击个人中心
    utils.Sleep(2000)
    tap(300, 400)  // 点击余额查询
    
    // 2. OCR识别余额数字
    screenshot := takeScreenshot()
    balance := extractBalanceFromOCR(screenshot)
    
    // 3. 可视化显示识别结果
    gui.HudSetText([]gui.TextItem{
        {TextColor: "#FFFFFF", Text: "当前余额: "},
        {TextColor: "#00FF00", Text: fmt.Sprintf("%.2f元", balance)},
    })
    
    // 4. 测试断言
    if balance >= 0 {
        gui.Toast("✅ 余额显示正确")
    } else {
        gui.Toast("❌ 余额显示异常")
    }
}
```

## 具体流程
// 1. 截图工具
screenshot := utils.Screenshot()

// 2. OCR引擎（前面的PPOCR）
textResults := ppocr.Detect(screenshot)

// 3. 坐标计算
clickX, clickY := calculateCenter(textResults[0].BoundingBox)

// 4. 触控模拟
utils.Tap(clickX, clickY)

// 5. 图形化反馈（使用这些GUI函数）
gui.DrawRect(clickX-10, clickY-10, clickX+10, clickY+10, "#FF0000")
gui.Toast(fmt.Sprintf("点击: %d,%d", clickX, clickY))

测试报告生成
func GenerateTestReport() {
    gui.HudSetText([]gui.TextItem{
        {TextColor: "#FFFFFF", Text: "测试报告: "},
        {TextColor: "#00FF00", Text: "通过 8 "},
        {TextColor: "#FF0000", Text: "失败 2 "},
        {TextColor: "#FFFF00", Text: "总计 10"},
    })
    
    // 在屏幕上绘制测试结果可视化
    gui.DrawRect(50, 500, 150, 550, "#00FF00")  // 通过
    gui.DrawRect(160, 500, 260, 550, "#FF0000") // 失败
}