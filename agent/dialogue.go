package agent

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// DialogueSystem 对话系统
type DialogueSystem struct {
	agent  *Agent
	memory *ConversationMemory
}

// NewDialogueSystem 创建对话系统
func NewDialogueSystem(agent *Agent) *DialogueSystem {
	return &DialogueSystem{
		agent:  agent,
		memory: NewConversationMemory(6),
	}
}

// Start 启动对话系统
func (ds *DialogueSystem) Start() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("AutoGo 自动化测试对话系统")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("输入 'exit' 或 'quit' 退出")
	fmt.Println("输入 'help' 查看帮助")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println()

	for {
		fmt.Print("有什么能帮您的吗？ ")
		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}

		// 处理特殊命令
		if query == "exit" || query == "quit" {
			fmt.Println("再见！")
			break
		}

		if query == "help" {
			ds.showHelp()
			continue
		}

		ds.memory.AddMessage("user", query)

		// 处理用户查询
		fmt.Println("\n正在处理您的请求...")
		result, err := ds.agent.ProcessQueryWithContext(query, ds.memory.ContextString())
		if err != nil {
			fmt.Printf("❌ 错误: %v\n\n", err)
			continue
		}

		// 显示结果
		fmt.Println(ds.agent.FormatResult(result))
		ds.memory.AddMessage("assistant", fmt.Sprintf("状态: %v, 报告: %s", result.Success, result.ReportPath))
		fmt.Println()
	}
}

// showHelp 显示帮助信息
func (ds *DialogueSystem) showHelp() {
	help := `
可用命令:
  - exit/quit: 退出系统
  - help: 显示此帮助信息

示例查询:
  - "点击登录按钮"
  - "在坐标(100, 200)点击"
  - "输入文本'Hello World'"
  - "验证页面是否存在'主页'文字"
  - "验证是否存在图片'logo.png'"
  - "启动应用com.example.app"
  - "等待5秒"
  - "滑动从(100,200)到(300,400)"

支持的模块:
  - motion: 触摸操作（点击、滑动等）
  - uiacc: UI控件识别和操作
  - opencv: 图像模板匹配
  - ppocr: OCR文字识别
  - app: 应用管理
  - ime: 输入法操作
`
	fmt.Println(help)
}

// ProcessSingleQuery 处理单个查询（用于API调用）
func (ds *DialogueSystem) ProcessSingleQuery(query string) (*TestResult, error) {
	ds.memory.AddMessage("user", query)
	result, err := ds.agent.ProcessQueryWithContext(query, ds.memory.ContextString())
	if err == nil {
		ds.memory.AddMessage("assistant", fmt.Sprintf("状态: %v", result.Success))
	}
	return result, err
}

