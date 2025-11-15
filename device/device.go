package device

import "C"
import (
	_ "embed"
	"fmt"
	"github.com/xiaocainiao633/Genie1.0--/rhino"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xiaocainiao633/Genie1.0--/utils"
)

var (
	Width         int    // 设备分辨率宽度，横屏和竖屏时的数值不同。
	Height        int    // 设备分辨率高度，横屏和竖屏时的数值不同。
	SdkInt        int    // 安卓系统API版本。例如安卓4.4的sdkInt为19
	CpuAbi        string // 设备的CPU架构，如"arm64-v8a", "x86", "x86_64"
	BuildId       string // 修订版本号，或者诸如"M4-rc20"的标识
	Broad         string // 设备的主板型号
	Brand         string // 与产品或硬件相关的厂商品牌，如"Xiaomi", "Huawei"等
	Device        string // 设备在工业设计中的名称
	Model         string // 设备型号
	Product       string // 整个产品的名称
	Bootloader    string // 设备Bootloader的版本
	Hardware      string // 设备的硬件名称
	Fingerprint   string // 构建(build)的唯一标识码
	Serial        string // 硬件序列号
	Incremental   string // 设备构建的内部版本号
	Release       string // Android系统版本号。例如"5.0", "7.1.1"
	BaseOS        string // 设备的基础操作系统版本
	SecurityPatch string // 安全补丁程序级别
	Codename      string // 开发代号，例如发行版是"REL"
)

func init() {
	Width, Height = utils.WmSize()
	props := utils.Shell("getprop")
	CpuAbi = utils.Shell("getprop ro.product.cpu.abi")
	BuildId = parseProp(props, "ro.build.id")
	Broad = parseProp(props, "ro.product.board")
	Brand = parseProp(props, "ro.product.brand")
	Device = parseProp(props, "ro.product.device")
	Model = parseProp(props, "ro.product.model")
	Product = parseProp(props, "ro.product.name")
	Bootloader = parseProp(props, "ro.bootloader")
	Hardware = parseProp(props, "ro.hardware")
	Fingerprint = parseProp(props, "ro.build.fingerprint")
	Serial = parseProp(props, "ro.serialno")
	SdkInt, _ = strconv.Atoi(parseProp(props, "ro.build.version.sdk"))
	Incremental = parseProp(props, "ro.build.version.incremental")
	Release = parseProp(props, "ro.build.version.release")
	BaseOS = parseProp(props, "ro.build.version.base_os")
	SecurityPatch = parseProp(props, "ro.build.version.security_patch")
	Codename = parseProp(props, "ro.build.version.codename")
	go heartbeat()
}

// 定时获取设备分辨率信息
func heartbeat() {
	for {
		Width, Height = utils.WmSize()
		time.Sleep(1 * time.Second)
	}
}

// GetImei 返回设备的IMEI。
func GetImei() string {
	arr := strings.Split(utils.Shell("service call iphonesubinfo 1"), "\n")
	if len(arr) < 3 {
		return ""
	}
	imei := ""
	for i := 1; i < 4; i++ {
		imei = imei + arr[i][51:]
	}
	imei = strings.Replace(imei, ".", "", -1)
	imei = strings.Replace(imei, "'", "", -1)
	imei = strings.Replace(imei, " ", "", -1)
	imei = strings.Replace(imei, ")", "", -1)
	return imei
}

// GetAndroidId 返回设备的Android ID。
func GetAndroidId() string {
	return utils.Shell("settings get secure android_id")
}

// GetWifiMac 获取设备WIFI-MAC
func GetWifiMac() string {
	name := []string{"wlan0", "wlan1", "p2p0", "eth0", "bond0", "dummy0"}
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, netName := range name {
		for _, inter := range interfaces {
			if inter.Name == netName {
				mac := inter.HardwareAddr.String()
				if len(mac) == 17 && mac != "00:00:00:00:00:00" {
					return mac
				}
			}
		}
	}
	return ""
}

// GetWlanMac 获取设备以太网MAC
func GetWlanMac() string {
	mac := utils.Shell("cat /sys/class/net/eth0/address")
	if len(mac) == 17 {
		return mac
	}
	return ""
}

// GetIp 获取设备局域网IP地址
func GetIp() string {
	// 使用 net.Interfaces() 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return ""
	}

	for _, iface := range interfaces {
		// 检查接口是否是UP状态，并且排除Loopback接口
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			// 获取接口的地址列表
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error getting addresses for interface:", iface.Name, err)
				continue
			}

			for _, addr := range addrs {
				// 转换为 *net.IPNet 类型并检查是否为 IPv4 地址
				ipNet, ok := addr.(*net.IPNet)
				if ok && ipNet.IP.To4() != nil {
					ip := ipNet.IP.String()
					// 检查是否是局域网IP地址
					if isPrivateIP(ipNet.IP) {
						return ip
					}
				}
			}
		}
	}

	return ""
}

// GetBrightness 返回当前的(手动)亮度。范围为0~255。
func GetBrightness() string {
	return utils.Shell("settings get system screen_brightness")
}

// GetBrightnessMode 返回当前亮度模式，0为手动亮度，1为自动亮度。
func GetBrightnessMode() string {
	return utils.Shell("settings get system screen_brightness_mode")
}

// GetMusicVolume 返回当前媒体音量。
func GetMusicVolume() int {
	return s2i(getMiddleString(utils.Shell("cmd media_session volume --stream 3 --get | grep range"), "volume is ", " "))
}

// GetNotificationVolume 返回当前通知音量。
func GetNotificationVolume() int {
	return s2i(getMiddleString(utils.Shell("cmd media_session volume --stream 5 --get | grep range"), "volume is ", " "))
}

// GetAlarmVolume 返回当前闹钟音量。
func GetAlarmVolume() int {
	return s2i(getMiddleString(utils.Shell("cmd media_session volume --stream 4 --get | grep range"), "volume is ", " "))
}

// GetMusicMaxVolume 返回媒体音量的最大值。
func GetMusicMaxVolume() int {
	return s2i(getMiddleString(utils.Shell("cmd media_session volume --stream 3 --get | grep range"), "..", "]"))
}

// GetNotificationMaxVolume 返回通知音量的最大值。
func GetNotificationMaxVolume() int {
	return s2i(getMiddleString(utils.Shell("cmd media_session volume --stream 5 --get | grep range"), "..", "]"))
}

// GetAlarmMaxVolume 返回闹钟音量的最大值。
func GetAlarmMaxVolume() int {
	return s2i(getMiddleString(utils.Shell("cmd media_session volume --stream 4 --get | grep range"), "..", "]"))
}

// SetMusicVolume 设置当前媒体音量。
func SetMusicVolume(volume int) {
	utils.Shell("cmd media_session volume --show --stream 3 --set " + i2s(volume))
}

// SetNotificationVolume 设置当前通知音量。
func SetNotificationVolume(volume int) {
	utils.Shell("cmd media_session volume --show --stream 5 --set " + i2s(volume))
}

// SetAlarmVolume 设置当前闹钟音量。
func SetAlarmVolume(volume int) {
	utils.Shell("cmd media_session volume --show --stream 4 --set " + i2s(volume))
}

// GetBattery 返回当前电量百分比。
func GetBattery() int {
	return s2i(getMiddleString(utils.Shell("dumpsys battery | grep level")+"|", ": ", "|"))
}

// GetBatteryStatus 获取电池状态。 1：没有充电；2：正充电；3：没插充电器；4：不充电； 5：电池充满
func GetBatteryStatus() int {
	return s2i(getMiddleString(utils.Shell("dumpsys battery | grep status")+"|", ": ", "|"))
}

// SetBatteryStatus 模拟电池状态。 1：没有充电；2：正充电；5：电池充满
func SetBatteryStatus(value int) {
	utils.Shell("dumpsys battery set status " + i2s(value))
}

// SetBatteryLevel 模拟电池电量百分百 0-100
func SetBatteryLevel(value int) {
	utils.Shell("dumpsys battery set level " + i2s(value))
}

// GetTotalMem 返回设备内存总量，单位(KB)。1MB = 1024KB。
func GetTotalMem() int {
	matches := regexp.MustCompile(`MemTotal:\s+(\d+)`).FindStringSubmatch(utils.Shell("cat /proc/meminfo"))
	if len(matches) > 1 {
		return s2i(matches[1])
	}
	return 0
}

// GetAvailMem 返回设备当前可用的内存，单位字节(KB)。
func GetAvailMem() int {
	matches := regexp.MustCompile(`MemAvailable:\s+(\d+)`).FindStringSubmatch(utils.Shell("cat /proc/meminfo"))
	if len(matches) > 1 {
		return s2i(matches[1])
	}
	return 0
}

// IsScreenOn 返回设备屏幕是否是亮着的。如果屏幕亮着，返回true; 否则返回false。
func IsScreenOn() bool {
	result := regexp.MustCompile(`=(.*)`).FindStringSubmatch(utils.Shell("dumpsys window policy | grep screenState="))
	if result[1] == "SCREEN_STATE_OFF" || result[1] == "0" {
		return false
	}
	return utils.Shell("dumpsys power | grep mWakefulness=Dozing") == ""
}

// IsScreenUnlock 返回屏幕锁是否已经解开。已经解开返回true; 否则返回false。
func IsScreenUnlock() bool {
	return utils.Shell("dumpsys window policy | grep showing=true") == ""
}

// SetScreenMode 设置屏幕显示模式
// mode取值范围：
// 0 = OFF(关闭屏幕)
// 1 = DOZE(低功耗待机亮屏模式)
// 2 = ON(正常点亮屏幕)
// 3 = DOZE_SUSPEND(在DOZE基础上进一步降低功耗的待机模式)
func SetScreenMode(mode int) {
	if mode < 0 || mode > 4 {
		return
	}
	rhino.Eval("_device", fmt.Sprintf(`
importClass(android.os.Build);
importClass(android.view.SurfaceControl);

var internalDisplayToken;
if (Build.VERSION.SDK_INT < 29) {
    internalDisplayToken = SurfaceControl.getBuiltInDisplay(0);
} else {
    internalDisplayToken = SurfaceControl.getInternalDisplayToken();
}
SurfaceControl.setDisplayPowerMode(internalDisplayToken, %d);
`, mode))
}

// WakeUp 唤醒设备，包括唤醒设备CPU、屏幕等，可以用来点亮屏幕。
func WakeUp() {
	utils.Shell("input keyevent KEYCODE_WAKEUP")
}

// KeepScreenOn 保持屏幕常亮。
func KeepScreenOn() {
	utils.Shell("settings put system screen_off_timeout 2147483647;svc power stayon true")
}

// Vibrate 使设备震动一段时间，单位毫秒，需要root权限。
func Vibrate(ms int) {
	utils.Shell("echo " + i2s(ms) + ">/sys/devices/virtual/timed_output/vibrator/enable")
}

// CancelVibration 如果设备处于震动状态，则取消震动，需要root权限。
func CancelVibration() {
	Vibrate(0)
}

// 判断是否为局域网IP地址
func isPrivateIP(ip net.IP) bool {
	privateIPBlocks := []*net.IPNet{
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// getMiddleString 获取字符串中间部分
func getMiddleString(str, starting, ending string) string {
	s := strings.Index(str, starting)
	if s < 0 {
		return ""
	}
	s += len(starting)
	e := strings.Index(str[s:], ending)
	if e < 0 {
		return ""
	}
	return str[s : s+e]
}

// parseProp 解析属性
func parseProp(props, key string) string {
	lines := strings.Split(props, "\n")
	for _, line := range lines {
		if strings.Contains(line, key) {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 {
				return strings.Trim(parts[1], "[]")
			}
		}
	}
	return ""
}

// s2i 将字符串转换为整数
func s2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

// i2s 将整数转换为字符串
func i2s(i int) string {
	return strconv.Itoa(i)
}
