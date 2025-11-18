package agent

import (
	"fmt"
	"strings"
)

// CodeGenerator 代码生成器
type CodeGenerator struct {
	kb     *KnowledgeBase
	ollama *OllamaClient
	useLLM bool
}

// NewCodeGenerator 创建代码生成器
func NewCodeGenerator(kb *KnowledgeBase) *CodeGenerator {
	return &CodeGenerator{kb: kb}
}

// EnableLLM 启用LLM生成
func (cg *CodeGenerator) EnableLLM(client *OllamaClient) {
	cg.ollama = client
	cg.useLLM = client != nil
}

// GenerateTestScript 根据用户输入生成测试脚本
func (cg *CodeGenerator) GenerateTestScript(userQuery string, memoryContext string) (string, error) {
	// 获取相关API文档
	context, err := cg.kb.GetContext(userQuery)
	if err != nil {
		return "", err
	}

	// 解析用户意图
	intent := cg.parseIntent(userQuery)
	
	// 生成代码
	if cg.useLLM && cg.ollama != nil {
		code, err := cg.generateCodeWithLLM(intent, context, memoryContext)
		if err == nil && strings.TrimSpace(code) != "" {
			return code, nil
		}
	}

	code := cg.generateCode(intent, context)

	return code, nil
}

// Intent 用户意图
type Intent struct {
	Action      string            // 操作类型: click, input, assert, wait, etc.
	Target      string            // 目标: button, text, image, etc.
	Value       string            // 值: 坐标、文本、图片路径等
	Module      string            // 使用的模块
	Parameters  map[string]string // 额外参数
}

func (i Intent) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Action: %s\n", i.Action))
	builder.WriteString(fmt.Sprintf("Target: %s\n", i.Target))
	if i.Value != "" {
		builder.WriteString(fmt.Sprintf("Value: %s\n", i.Value))
	}
	if i.Module != "" {
		builder.WriteString(fmt.Sprintf("Module: %s\n", i.Module))
	}
	if len(i.Parameters) > 0 {
		builder.WriteString("Parameters:\n")
		for k, v := range i.Parameters {
			builder.WriteString(fmt.Sprintf(" - %s: %s\n", k, v))
		}
	}
	return builder.String()
}

// parseIntent 解析用户意图
func (cg *CodeGenerator) parseIntent(query string) Intent {
	queryLower := strings.ToLower(query)
	intent := Intent{
		Parameters: make(map[string]string),
	}

	// 检测操作类型
	if strings.Contains(queryLower, "点击") || strings.Contains(queryLower, "click") {
		intent.Action = "click"
	} else if strings.Contains(queryLower, "输入") || strings.Contains(queryLower, "input") || strings.Contains(queryLower, "输入文本") {
		intent.Action = "input"
	} else if strings.Contains(queryLower, "验证") || strings.Contains(queryLower, "检查") || strings.Contains(queryLower, "assert") {
		intent.Action = "assert"
	} else if strings.Contains(queryLower, "等待") || strings.Contains(queryLower, "wait") {
		intent.Action = "wait"
	} else if strings.Contains(queryLower, "启动") || strings.Contains(queryLower, "launch") || strings.Contains(queryLower, "打开") {
		intent.Action = "launch"
	} else if strings.Contains(queryLower, "滑动") || strings.Contains(queryLower, "swipe") {
		intent.Action = "swipe"
	}

	// 检测目标类型
	if strings.Contains(queryLower, "按钮") || strings.Contains(queryLower, "button") {
		intent.Target = "button"
	} else if strings.Contains(queryLower, "输入框") || strings.Contains(queryLower, "input") || strings.Contains(queryLower, "文本框") {
		intent.Target = "input"
	} else if strings.Contains(queryLower, "图片") || strings.Contains(queryLower, "image") || strings.Contains(queryLower, "图像") {
		intent.Target = "image"
	} else if strings.Contains(queryLower, "文字") || strings.Contains(queryLower, "text") || strings.Contains(queryLower, "文本") {
		intent.Target = "text"
	}

	// 提取文本内容（引号内的内容）
	if idx := strings.Index(query, "\""); idx != -1 {
		if idx2 := strings.Index(query[idx+1:], "\""); idx2 != -1 {
			intent.Value = query[idx+1 : idx+idx2+1]
		}
	}

	// 提取坐标
	if strings.Contains(query, "(") && strings.Contains(query, ",") {
		// 尝试提取坐标
		parts := strings.Fields(query)
		for _, part := range parts {
			if strings.Contains(part, ",") && strings.Contains(part, "(") {
				intent.Value = part
				intent.Target = "coordinate"
				break
			}
		}
	}

	// 检测使用的模块
	if strings.Contains(queryLower, "ui") || strings.Contains(queryLower, "控件") {
		intent.Module = "uiacc"
	} else if strings.Contains(queryLower, "图像") || strings.Contains(queryLower, "图片") || strings.Contains(queryLower, "模板") {
		intent.Module = "opencv"
	} else if strings.Contains(queryLower, "ocr") || strings.Contains(queryLower, "识别") || strings.Contains(queryLower, "文字") {
		intent.Module = "ppocr"
	} else if strings.Contains(queryLower, "坐标") || strings.Contains(queryLower, "点击") {
		intent.Module = "motion"
	}

	return intent
}

// generateCode 生成代码
func (cg *CodeGenerator) generateCode(intent Intent, context string) string {
	var code strings.Builder

	// 添加包声明和导入
	code.WriteString("package main\n\n")
	code.WriteString("import (\n")
	code.WriteString("\t\"fmt\"\n")
	code.WriteString("\t\"os\"\n")
	code.WriteString("\t\"strings\"\n")
	code.WriteString("\t\"time\"\n\n")
	
	// 根据需要的模块添加导入
	modules := cg.getRequiredModules(intent)
	for _, module := range modules {
		code.WriteString(fmt.Sprintf("\t\"github.com/xiaocainiao633/Genie1.0--/tree/main/%s\"\n", module))
	}
	code.WriteString("\t\"github.com/xiaocainiao633/Genie1.0--/tree/main/utils\"\n")
	code.WriteString(")\n\n")

	// 添加主函数
	code.WriteString("func main() {\n")
	code.WriteString("\tfmt.Println(\"开始执行自动化测试...\")\n\n")

	// 根据意图生成代码
	switch intent.Action {
	case "click":
		code.WriteString(cg.generateClickCode(intent))
	case "input":
		code.WriteString(cg.generateInputCode(intent))
	case "assert":
		code.WriteString(cg.generateAssertCode(intent))
	case "wait":
		code.WriteString(cg.generateWaitCode(intent))
	case "launch":
		code.WriteString(cg.generateLaunchCode(intent))
	case "swipe":
		code.WriteString(cg.generateSwipeCode(intent))
	default:
		code.WriteString("\t// 未识别的操作类型\n")
	}

	code.WriteString("\tfmt.Println(\"测试完成\")\n")
	code.WriteString("}\n")

	return code.String()
}

func (cg *CodeGenerator) generateCodeWithLLM(intent Intent, context, memoryContext string) (string, error) {
	if cg.ollama == nil {
		return "", fmt.Errorf("LLM未配置")
	}

	var prompt strings.Builder
	prompt.WriteString("你是一名资深的Go语言自动化测试工程师。")
	prompt.WriteString("请根据以下信息生成一个完整的Go测试脚本，脚本会在AutoGo环境中执行。\n\n")
	if memoryContext != "" {
		prompt.WriteString(memoryContext + "\n")
	}
	prompt.WriteString("相关API文档:\n")
	prompt.WriteString(context + "\n")
	prompt.WriteString("用户需求:\n")
	prompt.WriteString(intent.String() + "\n\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("1. 必须包含package main和main函数。\n")
	prompt.WriteString("2. 导入必要的AutoGo模块。\n")
	prompt.WriteString("3. 代码可直接编译运行。\n")
	prompt.WriteString("4. 添加必要的错误处理和日志输出。\n")

	response, err := cg.ollama.Generate(prompt.String())
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

// getRequiredModules 获取需要的模块
func (cg *CodeGenerator) getRequiredModules(intent Intent) []string {
	modules := make(map[string]bool)
	
	if intent.Module != "" {
		modules[intent.Module] = true
	}
	
	// 根据操作类型添加必要模块
	switch intent.Action {
	case "click":
		if intent.Target == "image" {
			modules["opencv"] = true
			modules["images"] = true
		}
		modules["motion"] = true
	case "input":
		modules["uiacc"] = true
		modules["ime"] = true
	case "assert":
		if intent.Target == "text" {
			modules["ppocr"] = true
			modules["images"] = true
		} else if intent.Target == "image" {
			modules["opencv"] = true
			modules["images"] = true
		}
	case "launch":
		modules["app"] = true
	}

	result := make([]string, 0, len(modules))
	for module := range modules {
		result = append(result, module)
	}
	return result
}

// generateClickCode 生成点击代码
func (cg *CodeGenerator) generateClickCode(intent Intent) string {
	var code strings.Builder

	if intent.Target == "coordinate" {
		// 直接坐标点击
		code.WriteString("\t// 在指定坐标点击\n")
		code.WriteString(fmt.Sprintf("\tmotion.Click(%s, 1)\n", intent.Value))
	} else if intent.Target == "button" && intent.Value != "" {
		// 通过文本查找按钮
		code.WriteString("\t// 查找并点击按钮\n")
		code.WriteString(fmt.Sprintf("\tobj := uiacc.New().Text(\"%s\").FindOnce()\n", intent.Value))
		code.WriteString("\tif obj != nil {\n")
		code.WriteString("\t\tobj.Click()\n")
		code.WriteString("\t\tfmt.Println(\"点击成功\")\n")
		code.WriteString("\t} else {\n")
		code.WriteString("\t\tfmt.Println(\"未找到按钮\")\n")
		code.WriteString("\t}\n")
	} else if intent.Target == "image" && intent.Value != "" {
		// 通过图像匹配点击
		code.WriteString("\t// 通过图像匹配查找并点击\n")
		code.WriteString(fmt.Sprintf("\ttemplateBytes, err := os.ReadFile(\"%s\")\n", intent.Value))
		code.WriteString("\tif err != nil {\n")
		code.WriteString("\t\tfmt.Printf(\"读取模板失败: %%v\\n\", err)\n")
		code.WriteString("\t\treturn\n")
		code.WriteString("\t}\n")
		code.WriteString("\tx, y := opencv.FindImage(0, 0, 0, 0, &templateBytes, false, 1.0, 0.8)\n")
		code.WriteString("\tif x != -1 && y != -1 {\n")
		code.WriteString("\t\tmotion.Click(x, y, 1)\n")
		code.WriteString("\t\tfmt.Printf(\"找到并点击，坐标: (%%d, %%d)\\n\", x, y)\n")
		code.WriteString("\t} else {\n")
		code.WriteString("\t\tfmt.Println(\"未找到目标图像\")\n")
		code.WriteString("\t}\n")
	} else {
		code.WriteString("\t// 通用点击代码\n")
		code.WriteString("\t// 请根据实际情况修改\n")
	}

	return code.String()
}

// generateInputCode 生成输入代码
func (cg *CodeGenerator) generateInputCode(intent Intent) string {
	var code strings.Builder

	if intent.Value != "" {
		code.WriteString("\t// 输入文本\n")
		if intent.Target == "input" {
			code.WriteString("\t// 查找输入框\n")
			code.WriteString("\tinputObj := uiacc.New().Editable(true).FindOnce()\n")
			code.WriteString("\tif inputObj != nil {\n")
			code.WriteString(fmt.Sprintf("\t\tinputObj.SetText(\"%s\")\n", intent.Value))
			code.WriteString("\t\tfmt.Println(\"输入成功\")\n")
			code.WriteString("\t} else {\n")
			code.WriteString("\t\tfmt.Println(\"未找到输入框\")\n")
			code.WriteString("\t}\n")
		} else {
			code.WriteString(fmt.Sprintf("\time.InputText(\"%s\")\n", intent.Value))
		}
	} else {
		code.WriteString("\t// 输入代码（需要指定文本内容）\n")
	}

	return code.String()
}

// generateAssertCode 生成断言代码
func (cg *CodeGenerator) generateAssertCode(intent Intent) string {
	var code strings.Builder

	if intent.Target == "text" && intent.Value != "" {
		code.WriteString("\t// 验证文本是否存在\n")
		code.WriteString("\timg := images.CaptureScreen(0, 0, 0, 0)\n")
		code.WriteString("\tresults := ppocr.OcrFromImage(img, \"\")\n")
		code.WriteString("\tfound := false\n")
		code.WriteString("\tfor _, result := range results {\n")
		code.WriteString(fmt.Sprintf("\t\tif strings.Contains(result.Label, \"%s\") {\n", intent.Value))
		code.WriteString("\t\t\tfound = true\n")
		code.WriteString("\t\t\tfmt.Printf(\"找到文本: %%s\\n\", result.Label)\n")
		code.WriteString("\t\t\tbreak\n")
		code.WriteString("\t\t}\n")
		code.WriteString("\t}\n")
		code.WriteString("\tif found {\n")
		code.WriteString("\t\tfmt.Println(\"✅ 断言通过：文本存在\")\n")
		code.WriteString("\t} else {\n")
		code.WriteString("\t\tfmt.Println(\"❌ 断言失败：文本不存在\")\n")
		code.WriteString("\t}\n")
	} else if intent.Target == "image" && intent.Value != "" {
		code.WriteString("\t// 验证图像是否存在\n")
		code.WriteString(fmt.Sprintf("\ttemplateBytes, err := os.ReadFile(\"%s\")\n", intent.Value))
		code.WriteString("\tif err != nil {\n")
		code.WriteString("\t\tfmt.Printf(\"读取模板失败: %%v\\n\", err)\n")
		code.WriteString("\t\treturn\n")
		code.WriteString("\t}\n")
		code.WriteString("\tx, y := opencv.FindImage(0, 0, 0, 0, &templateBytes, false, 1.0, 0.8)\n")
		code.WriteString("\tif x != -1 && y != -1 {\n")
		code.WriteString("\t\tfmt.Printf(\"✅ 断言通过：图像存在，坐标: (%%d, %%d)\\n\", x, y)\n")
		code.WriteString("\t} else {\n")
		code.WriteString("\t\tfmt.Println(\"❌ 断言失败：图像不存在\")\n")
		code.WriteString("\t}\n")
	} else {
		code.WriteString("\t// 断言代码（需要指定验证内容）\n")
	}

	return code.String()
}

// generateWaitCode 生成等待代码
func (cg *CodeGenerator) generateWaitCode(intent Intent) string {
	var code strings.Builder

	timeout := "5000"
	if t, ok := intent.Parameters["timeout"]; ok {
		timeout = t
	}

	if intent.Value != "" {
		code.WriteString("\t// 等待元素出现\n")
		code.WriteString(fmt.Sprintf("\tobj := uiacc.New().Text(\"%s\").WaitFor(%s)\n", intent.Value, timeout))
		code.WriteString("\tif obj != nil {\n")
		code.WriteString("\t\tfmt.Println(\"元素已出现\")\n")
		code.WriteString("\t} else {\n")
		code.WriteString("\t\tfmt.Println(\"等待超时\")\n")
		code.WriteString("\t}\n")
	} else {
		code.WriteString(fmt.Sprintf("\t// 等待 %s 毫秒\n", timeout))
		code.WriteString(fmt.Sprintf("\tutils.Sleep(%s)\n", timeout))
	}

	return code.String()
}

// generateLaunchCode 生成启动应用代码
func (cg *CodeGenerator) generateLaunchCode(intent Intent) string {
	var code strings.Builder

	if intent.Value != "" {
		code.WriteString("\t// 启动应用\n")
		code.WriteString(fmt.Sprintf("\tif app.Launch(\"%s\", 0) {\n", intent.Value))
		code.WriteString("\t\tfmt.Println(\"应用启动成功\")\n")
		code.WriteString("\t\tutils.Sleep(2000)\n")
		code.WriteString("\t} else {\n")
		code.WriteString("\t\tfmt.Println(\"应用启动失败\")\n")
		code.WriteString("\t\treturn\n")
		code.WriteString("\t}\n")
	} else {
		code.WriteString("\t// 启动应用（需要指定包名）\n")
	}

	return code.String()
}

// generateSwipeCode 生成滑动代码
func (cg *CodeGenerator) generateSwipeCode(intent Intent) string {
	var code strings.Builder

	code.WriteString("\t// 执行滑动操作\n")
	code.WriteString("\t// 请根据实际情况修改坐标\n")
	code.WriteString("\tmotion.Swipe(100, 200, 300, 400, 500)\n")

	return code.String()
}

