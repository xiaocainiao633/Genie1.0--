package system

// 提供与系统相关的功能，包括进程管理和资源查询

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/xiaocainiao633/Genie1.0--/files"
	"github.com/xiaocainiao633/Genie1.0--/utils"
)

// 获取指定进程的 PID
func GetPid(processName string) int {
	if processName == "" {
		return os.Getpid()
	}
	var err error
	pid := 0
	output := utils.Shell("ps -ef")
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasSuffix(line, processName) {
			fields := strings.Fields(line)
			pid, err = strconv.Atoi(fields[1])
			if err == nil {
				return pid
			}
			break
		}
	}

	output = utils.Shell("ps")
	lines = strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasSuffix(line, processName) {
			fields := strings.Fields(line)
			pid, err = strconv.Atoi(fields[1])
			if err == nil {
				return pid
			}
			return -1
		}
	}
	return -1
}

// 获取指定进程的内存使用量
func GetMemoryUsage(pid int) int {
	if pid == 0 {
		pid = os.Getpid()
	}

	cmd := fmt.Sprintf("cat /proc/%d/status | grep -e VmRSS", pid)
	output := utils.Shell(cmd)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.Contains(line, "VmRSS") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memory, err := strconv.Atoi(fields[1])
				if err == nil {
					return memory
				}
			}
		}
	}

	return -1
}

// 获取指定进程的 CPU 使用率
func GetCpuUsage(pid int) float64 {
	if pid == 0 {
		pid = os.Getpid()
	}

	cmd := fmt.Sprintf("top -b -n 1 | grep '^ *%d '", pid)
	output := utils.Shell(cmd)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 8 {
			cpuUsage, err := strconv.ParseFloat(fields[8], 64)
			if err == nil {
				return cpuUsage
			}
		}
	}

	return 0.0 // 返回 0.0 表示查询失败
}

// 重启当前脚本进程
func RestartSelf() {
	os.Exit(123)
}

// 设置脚本开机自动运行
func SetBootStart(enable bool) {
	if !enable {
		utils.Shell("rm -rf /data/local/tmp/start")
		return
	}

	if !strings.Contains(utils.Shell("whoami"), "root") {
		fmt.Println("开机自启设置失败,无root权限")
		return
	}

	var cmd = "#!/system/bin/sh\nwhile [ \"$(getprop sys.boot_completed)\" != \"1\" ]; do\n    sleep 1\ndone\nsleep 10\n"

	if os.Getenv("APPPID") != "" {
		dir := filepath.Dir(os.Args[0])
		pkg := filepath.Base(dir)
		cmd = cmd + "sh /sdcard/Android/data/" + pkg + "/files/run.sh\nam start -n $(cmd package resolve-activity --brief " + pkg + " android.intent.action.MAIN | grep " + pkg + "/) --ei \"run\" 1"
	} else {
		cmd = cmd + os.Args[0] + " -bootstart=1"
	}

	_ = os.WriteFile("/data/local/tmp/start", []byte(cmd), 0644)

	if dirExists("/data/adb/modules") {
		utils.Shell("mkdir /data/adb/service.d")
		_ = os.WriteFile("/data/adb/service.d/AutoGo.sh", []byte("sh /data/local/tmp/start"), 0644)
		utils.Shell("chmod 755 /data/adb/service.d/AutoGo.sh")
		return
	}

	cil := `
	
	(typepermissive adbd)
	(typepermissive shell)
	
	;; +exec typeattributeset file_type, exec_type, mlstrustedobject
	;; + typeattributeset domain, mlstrustedobject, mlstrustedsubject, netdomain, coredomain
	
	`

	rc := `

on property:sys.boot_completed=1
    start autogo
    
service autogo /system/bin/sh /data/local/tmp/start
    seclabel u:r:su:s0
    oneshot
	
`
	if strings.Contains(files.Read("/system/etc/selinux/plat_sepolicy.cil"), cil) && strings.Contains(files.Read("/system/etc/init/logd.rc"), rc) {
		return
	}

	utils.Shell("mount -o rw,remount /system")
	str := utils.Shell("mount -o rw,remount /")
	if str != "" && (runtime.GOARCH == "amd64" || runtime.GOARCH == "386") {
		fmt.Println("系统分区挂载失败,请在模拟器性能设置中开启System.vmdk可写入")
		return
	}

	if files.Exists("/system/etc/selinux/plat_sepolicy.cil") && !strings.Contains(files.Read("/system/etc/selinux/plat_sepolicy.cil"), cil) {
		utils.Shell("rm -rf /system/etc/selinux/*.sha256")
		files.Append("/system/etc/selinux/plat_sepolicy.cil", cil)
	}

	if !strings.Contains(files.Read("/system/etc/init/logd.rc"), rc) {
		files.Append("/system/etc/init/logd.rc", rc)
	}

	utils.Shell("mount -o ro,remount /system")
	utils.Shell("mount -o ro,remount /")
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// 如果出错，可能是不存在，也可能是权限问题
		return false
	}
	return info.IsDir()
}
