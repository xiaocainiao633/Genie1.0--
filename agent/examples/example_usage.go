package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/xiaocainiao633/Genie1.0--/agent"
)

func main() {
	// 1. 初始化知识库（首次运行）
	kbPath := "./knowledge_base.db"
	workspace := "./workspace"

	// 检查知识库是否存在，不存在则初始化
	if _, err := os.Stat(kbPath); os.IsNotExist(err) {
		fmt.Println("初始化知识库...")
		kb, err := agent.NewKnowledgeBase(kbPath)
		if err != nil {
			fmt.Printf("创建知识库失败: %v\n", err)
			return
		}
		defer kb.Close()

		if err := agent.BuildDefaultKnowledgeBase(kb); err != nil {
			fmt.Printf("构建知识库失败: %v\n", err)
			return
		}
		fmt.Println("知识库初始化完成！")
	}

	// 2. 创建Agent
	ag, err := agent.NewAgent(kbPath, workspace)
	if err != nil {
		fmt.Printf("创建Agent失败: %v\n", err)
		return
	}
	defer ag.Close()

	// 3. 示例查询
	queries := []string{
		"点击登录按钮",
		"在输入框中输入'Hello World'",
		"验证页面是否存在'主页'文字",
		"启动应用com.example.app",
		"点击图片button.png",
	}

	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println("AutoGo RAG 系统示例")
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Println()

	for i, query := range queries {
		fmt.Printf("示例 %d: %s\n", i+1, query)
		fmt.Println("-" + strings.Repeat("-", 60))

		result, err := ag.ProcessQuery(query)
		if err != nil {
			fmt.Printf("❌ 错误: %v\n\n", err)
			continue
		}

		fmt.Println(ag.FormatResult(result))
		fmt.Println()
	}
}

