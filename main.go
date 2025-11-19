package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xiaocainiao633/Genie1.0--/agent"
)

func main() {
	var (
		kbPath       = flag.String("kb", "./knowledge_base.db", "知识库路径")
		workspace    = flag.String("workspace", "./workspace", "工作目录")
		reportDir    = flag.String("reports", "./workspace/reports", "报告目录")
		query        = flag.String("query", "", "单次查询（非交互模式）")
		initKB       = flag.Bool("init", false, "初始化知识库")
		useLLM       = flag.Bool("use-llm", true, "启用本地LLM辅助生成")
		autoExec     = flag.Bool("auto-exec", false, "自动在连接的Android设备上执行")
		adbPath      = flag.String("adb", "adb", "ADB命令路径")
		remoteDir    = flag.String("remote", "/data/local/tmp", "设备端执行目录")
		ollamaBase   = flag.String("ollama-base", "http://localhost:11434", "Ollama服务地址")
		ollamaModel  = flag.String("ollama-model", "llama3.2:latest", "Ollama推理模型")
		ollamaEmbed  = flag.String("ollama-embed", "llama3.2:latest", "Ollama向量模型")
	)
	flag.Parse()

	if err := os.MkdirAll(*workspace, 0755); err != nil {
		fmt.Printf("创建工作目录失败: %v\n", err)
		os.Exit(1)
	}

	if *initKB {
		initKnowledgeBase(*kbPath)
		return
	}

	var ollamaClient *agent.OllamaClient
	if *useLLM || *autoExec {
		ollamaClient = agent.NewOllamaClient(*ollamaBase, *ollamaModel, *ollamaEmbed)
	}

	cfg := agent.AgentConfig{
		UseLLM:       *useLLM,
		AutoExecute:  *autoExec,
		WorkspaceDir: *workspace,
		ReportDir:    *reportDir,
		ADBPath:      *adbPath,
		RemoteDir:    *remoteDir,
	}

	ag, err := agent.NewAgentWithOptions(*kbPath, cfg, ollamaClient)
	if err != nil {
		fmt.Printf("创建Agent失败: %v\n", err)
		os.Exit(1)
	}
	defer ag.Close()

	if *query != "" {
		result, err := ag.ProcessQueryWithContext(*query, "")
		if err != nil {
			fmt.Printf("处理查询失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(ag.FormatResult(result))
		return
	}

	dialogue := agent.NewDialogueSystem(ag)
	dialogue.Start()
}

func initKnowledgeBase(path string) {
	fmt.Println("正在初始化知识库...")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		os.Exit(1)
	}

	kb, err := agent.NewKnowledgeBase(path)
	if err != nil {
		fmt.Printf("创建知识库失败: %v\n", err)
		os.Exit(1)
	}
	defer kb.Close()

	if err := agent.BuildDefaultKnowledgeBase(kb); err != nil {
		fmt.Printf("构建知识库失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("知识库初始化完成！")
}