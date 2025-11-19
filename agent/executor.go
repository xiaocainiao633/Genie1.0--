package agent

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

// AndroidExecutor 负责将测试二进制推送到 Android 设备执行
type AndroidExecutor struct {
	ADBPath    string
	RemoteDir  string
	PackageDir string
}

// NewAndroidExecutor 创建执行器
func NewAndroidExecutor(adbPath, remoteDir string) *AndroidExecutor {
	if adbPath == "" {
		adbPath = "adb"
	}
	if remoteDir == "" {
		remoteDir = "/data/local/tmp"
	}
	return &AndroidExecutor{
		ADBPath:   adbPath,
		RemoteDir: remoteDir,
	}
}

// Execute 在设备上运行二进制
func (e *AndroidExecutor) Execute(localBinary string) (string, error) {
	remoteBinary := filepath.ToSlash(filepath.Join(e.RemoteDir, filepath.Base(localBinary)))

	pushCmd := exec.Command(e.ADBPath, "push", localBinary, remoteBinary)
	if output, err := pushCmd.CombinedOutput(); err != nil {
		return string(output), fmt.Errorf("adb push失败: %w", err)
	}

	chmodCmd := exec.Command(e.ADBPath, "shell", "chmod", "+x", remoteBinary)
	if output, err := chmodCmd.CombinedOutput(); err != nil {
		return string(output), fmt.Errorf("adb chmod失败: %w", err)
	}

	runCmd := exec.Command(e.ADBPath, "shell", remoteBinary)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	if err := runCmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("adb shell执行失败: %w", err)
	}

	return stdout.String(), nil
}

