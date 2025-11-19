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
	if usernameInput != nil {
		usernameInput.SetText("testuser")
	}
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

