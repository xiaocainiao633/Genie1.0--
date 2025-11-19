package agent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Agent 自动化测试Agent
type Agent struct {
	kb        *KnowledgeBase
	codeGen   *CodeGenerator
	options   AgentConfig
	executor  *AndroidExecutor
	reportDir string
}

// NewAgent 创建Agent（兼容旧接口）
func NewAgent(kbPath, workspaceDir string) (*Agent, error) {
	cfg := AgentConfig{WorkspaceDir: workspaceDir}
	cfg.normalize()
	return NewAgentWithOptions(kbPath, cfg, nil)
}

// NewAgentWithOptions 使用配置创建Agent
func NewAgentWithOptions(kbPath string, cfg AgentConfig, ollama *OllamaClient) (*Agent, error) {
	cfg.normalize()
	kb, err := NewKnowledgeBase(kbPath)
	if err != nil {
		return nil, err
	}

	if ollama != nil {
		kb.SetEmbedder(ollama)
	}

	// 初始化知识库
	docs, _ := kb.Search("click", 1)
	if len(docs) == 0 {
		if err := BuildDefaultKnowledgeBase(kb); err != nil {
			return nil, fmt.Errorf("构建知识库失败: %v", err)
		}
	}

	codeGen := NewCodeGenerator(kb)
	if cfg.UseLLM && ollama != nil {
		codeGen.EnableLLM(ollama)
	}

	var executor *AndroidExecutor
	if cfg.AutoExecute {
		executor = NewAndroidExecutor(cfg.ADBPath, cfg.RemoteDir)
	}

	return &Agent{
		kb:        kb,
		codeGen:   codeGen,
		options:   cfg,
		executor:  executor,
		reportDir: cfg.ReportDir,
	}, nil
}

// TestResult 测试结果
type TestResult struct {
	Success   bool          `json:"success"`
	Code      string        `json:"code"`
	Output    string        `json:"output"`
	Error     string        `json:"error,omitempty"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	ReportPath string       `json:"report_path,omitempty"`
}

// ProcessQuery 处理用户查询
func (a *Agent) ProcessQuery(userQuery string) (*TestResult, error) {
	return a.ProcessQueryWithContext(userQuery, "")
}

// ProcessQueryWithContext 带上下文处理
func (a *Agent) ProcessQueryWithContext(userQuery, memoryContext string) (*TestResult, error) {
	startTime := time.Now()

	// 1. 生成测试代码
	fmt.Println("[Agent] 正在分析用户需求...")
	code, err := a.codeGen.GenerateTestScript(userQuery, memoryContext)
	if err != nil {
		return &TestResult{
			Success:   false,
			Code:      "",
			Error:     fmt.Sprintf("代码生成失败: %v", err),
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, nil
	}

	fmt.Println("[Agent] 代码生成完成")
	fmt.Println("生成的代码:")
	fmt.Println("---")
	fmt.Println(code)
	fmt.Println("---")

	// 2. 保存代码到文件
	filename := fmt.Sprintf("generated_%d.go", time.Now().UnixNano())
	testFile := filepath.Join(a.options.WorkspaceDir, filename)
	if err := os.MkdirAll(a.options.WorkspaceDir, 0755); err != nil {
		return &TestResult{
			Success:   false,
			Code:      code,
			Error:     fmt.Sprintf("创建工作目录失败: %v", err),
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, nil
	}
	if err := os.WriteFile(testFile, []byte(code), 0644); err != nil {
		return &TestResult{
			Success:   false,
			Code:      code,
			Error:     fmt.Sprintf("保存代码失败: %v", err),
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, nil
	}

	// 3. 编译代码
	fmt.Println("[Agent] 正在编译代码...")
	binaryPath := filepath.Join(a.options.WorkspaceDir, fmt.Sprintf("test_binary_%d", time.Now().UnixNano()))
	buildCmd := exec.Command("go", "build", "-o", binaryPath, testFile)
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		return &TestResult{
			Success:   false,
			Code:      code,
			Error:     fmt.Sprintf("编译失败: %v\n输出: %s", err, string(buildOutput)),
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, nil
	}

	fmt.Println("[Agent] 编译成功")

	// 4. 执行测试（可选，在实际Android设备上运行）
	// 这里我们只返回生成的代码，实际执行需要部署到设备
	execOutput := "代码生成并编译成功。"
	var execErr error

	if a.executor != nil && a.options.AutoExecute {
		fmt.Println("[Agent] 正在推送到设备执行...")
		output, runErr := a.executor.Execute(binaryPath)
		execOutput = output
		execErr = runErr
	}

	success := execErr == nil
	if !a.options.AutoExecute {
		success = true
	}

	report := &TestReport{
		Query:           userQuery,
		CodePath:        testFile,
		BinaryPath:      binaryPath,
		CompileOutput:   string(buildOutput),
		ExecutionOutput: execOutput,
		ExecutionError:  errString(execErr),
		Success:         success,
		Duration:        time.Since(startTime),
		Timestamp:       time.Now(),
		AutoExecuted:    a.options.AutoExecute,
	}

	reportPath, _ := report.Save(a.reportDir)

	result := &TestResult{
		Success:    success,
		Code:       code,
		Output:     execOutput,
		Error:      errString(execErr),
		Duration:   time.Since(startTime),
		Timestamp:  report.Timestamp,
		ReportPath: reportPath,
	}

	return result, nil
}

// ProcessQueryWithExecution 处理查询并执行（如果可能）
func (a *Agent) ProcessQueryWithExecution(userQuery string) (*TestResult, error) {
	return a.ProcessQueryWithContext(userQuery, "")
}

// GetAPIInfo 获取API信息
func (a *Agent) GetAPIInfo(query string) ([]APIDoc, error) {
	return a.kb.Search(query, 10)
}

// Close 关闭Agent
func (a *Agent) Close() error {
	return a.kb.Close()
}

// FormatResult 格式化测试结果
func (a *Agent) FormatResult(result *TestResult) string {
	var output strings.Builder

	output.WriteString("=" + strings.Repeat("=", 60) + "\n")
	output.WriteString("自动化测试结果\n")
	output.WriteString("=" + strings.Repeat("=", 60) + "\n\n")

	if result.Success {
		output.WriteString("状态: ✅ 成功\n")
	} else {
		output.WriteString("状态: ❌ 失败\n")
	}

	output.WriteString(fmt.Sprintf("耗时: %v\n", result.Duration))
	output.WriteString(fmt.Sprintf("时间: %s\n\n", result.Timestamp.Format("2006-01-02 15:04:05")))

	output.WriteString("生成的代码:\n")
	output.WriteString("-" + strings.Repeat("-", 60) + "\n")
	output.WriteString(result.Code)
	output.WriteString("\n" + strings.Repeat("-", 60) + "\n\n")

	if result.Output != "" {
		output.WriteString("输出:\n")
		output.WriteString(result.Output + "\n\n")
	}

	if result.Error != "" {
		output.WriteString("错误:\n")
		output.WriteString(result.Error + "\n\n")
	}

	if result.ReportPath != "" {
		output.WriteString(fmt.Sprintf("报告: %s\n\n", result.ReportPath))
	}

	output.WriteString("=" + strings.Repeat("=", 60) + "\n")

	return output.String()
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

