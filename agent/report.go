package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestReport 自动化测试报告
type TestReport struct {
	Query            string        `json:"query"`
	CodePath         string        `json:"code_path"`
	BinaryPath       string        `json:"binary_path"`
	CompileOutput    string        `json:"compile_output"`
	ExecutionOutput  string        `json:"execution_output"`
	ExecutionError   string        `json:"execution_error,omitempty"`
	Success          bool          `json:"success"`
	Duration         time.Duration `json:"duration"`
	Timestamp        time.Time     `json:"timestamp"`
	AutoExecuted     bool          `json:"auto_executed"`
	AndroidDeviceLog string        `json:"android_device_log,omitempty"`
}

// Save 保存报告
func (r *TestReport) Save(reportDir string) (string, error) {
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("report-%s.json", r.Timestamp.Format("20060102-150405"))
	path := filepath.Join(reportDir, filename)

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
}

// Markdown 返回Markdown格式报告
func (r *TestReport) Markdown() string {
	var builder strings.Builder
	builder.WriteString("# 测试报告\n\n")
	builder.WriteString(fmt.Sprintf("- 时间: %s\n", r.Timestamp.Format("2006-01-02 15:04:05")))
	builder.WriteString(fmt.Sprintf("- 查询: %s\n", r.Query))
	builder.WriteString(fmt.Sprintf("- 耗时: %v\n", r.Duration))
	builder.WriteString(fmt.Sprintf("- 结果: %s\n\n", map[bool]string{true: "✅ 成功", false: "❌ 失败"}[r.Success]))

	if r.CodePath != "" {
		builder.WriteString(fmt.Sprintf("**代码文件**: `%s`\n\n", r.CodePath))
	}
	if r.BinaryPath != "" {
		builder.WriteString(fmt.Sprintf("**二进制文件**: `%s`\n\n", r.BinaryPath))
	}
	if r.CompileOutput != "" {
		builder.WriteString("## 编译输出\n```\n" + r.CompileOutput + "\n```\n\n")
	}
	if r.ExecutionOutput != "" {
		builder.WriteString("## 执行输出\n```\n" + r.ExecutionOutput + "\n```\n\n")
	}
	if r.ExecutionError != "" {
		builder.WriteString("## 错误信息\n```\n" + r.ExecutionError + "\n```\n\n")
	}

	return builder.String()
}

