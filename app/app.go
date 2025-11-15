package app

import (
	"mime"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/xiaocainiao633/Genie1.0--/utils"
)

type IntentOptions struct {
	Action      string
	Type        string
	Data        string
	Category    []string
	PackageName string
	ClassName   string
	Extras      map[string]string
	Flags       []string
}

// CurrentPackage 获取当前页面应用包名
func CurrentPackage() string {
	re := regexp.MustCompile(`mCurrentFocus=Window\{[^}]+\s([^\s/]+)/([^\s}]+)`)
	output := utils.Shell("dumpsys window | grep mCurrentFocus")
	match := re.FindStringSubmatch(output)
	if len(match) > 2 {
		return match[1]
	}
	return ""
}

// CurrentActivity 获取当前页面应用类名
func CurrentActivity() string {
	re := regexp.MustCompile(`mCurrentFocus=Window\{[^}]+\s([^\s/]+)/([^\s}]+)`)
	output := utils.Shell("dumpsys window | grep mCurrentFocus")
	match := re.FindStringSubmatch(output)
	if len(match) > 2 {
		return match[2]
	}
	return ""
}

// Launch 通过应用包名在指定页面启动应用。如果该包名对应的应用不存在，则返回false；否则返回true。
func Launch(packageName string, displayId int) bool {
	// 构建解析Activity的命令
	resolveCmd := "cmd package resolve-activity --brief " + packageName + " android.intent.action.MAIN"
	
	// 获取Activity路径
	activityPath := utils.Shell(resolveCmd + " | grep " + packageName)
	if activityPath == "" {
			return false
	}
	
	// 构建启动命令，添加displayId参数
	launchCmd := fmt.Sprintf("am start -n %s --display %d", activityPath, displayId)
	output := utils.Shell(launchCmd)
	
	return strings.Contains(output, "Starting")
}

// OpenAppSetting 打开应用的详情页(设置页)。如果找不到该应用，返回false; 否则返回true。
func OpenAppSetting(packageName string) bool {
	return !strings.Contains(utils.Shell(`am start -a android.settings.APPLICATION_DETAILS_SETTINGS -d "package:`+packageName+`"`), "Error:")
}

// ViewFile 用其他应用查看文件。文件不存在的情况由查看文件的应用处理。
func ViewFile(path string) {
	StartActivity(IntentOptions{
		Action: "VIEW",
		Type:   getMimeType(path),
		Data:   "file://" + path,
	})
}

// EditFile 用其他应用编辑文件。文件不存在的情况由编辑文件的应用处理
func EditFile(path string) {
	StartActivity(IntentOptions{
		Action: "EDIT",
		Type:   getMimeType(path),
		Data:   "file://" + path,
	})
}

// Uninstall 卸载应用
func Uninstall(packageName string) {
	utils.Shell("pm uninstall " + packageName)
}

// IsUninstalled 判断应用是否已卸载
func IsUninstalled(packageName string) bool {
	return utils.Shell("pm list packages | grep " + packageName) == ""
}

// Install 安装应用
func Install(path string) {
	utils.Shell("cp -rf " + path + " /data/local/tmp/1.apk;pm install -r /data/local/tmp/1.apk")
}

// IsInstalled 判断是否已经安装某个应用
func IsInstalled(packageName string) bool {
	return utils.Shell("pm list packages | grep "+packageName) != ""
}

// Clear 清除应用数据
func Clear(packageName string) {
	utils.Shell("pm clear " + packageName)
}

// IsCleared 判断是否清除应用数据
func IsCleared(packageName string) bool {
	return utils.Shell("pm clear " + packageName) == ""
}

// ForceStop 强制停止应用
func ForceStop(packageName string) {
	utils.Shell("am force-stop " + packageName)
}

// IsForceStopped 判断应用是否被强制停止
func IsForceStopped(packageName string) bool {
	return utils.Shell("am force-stop " + packageName) == ""
}

// Disable 禁用应用
func Disable(packageName string) {
	utils.Shell("pm disable-user " + packageName)
}

// IsDisabled 判断应用是否被禁用
func IsDisabled(packageName string) bool {
	// 是否能正常启动，不能启动则返回true
	return Launch(packageName, 0) == false
}

// IgnoreBattOpt 忽略电池优化
func IgnoreBattOpt(packageName string) {
	utils.Shell("dumpsys deviceidle whitelist +" + packageName)
}

// CancelIgnoreBattOpt 取消忽略电池优化
func CancelIgnoreBattOpt(packageName string) {
	utils.Shell("dumpsys deviceidle whitelist -" + packageName)
}

// Enable 启用应用，将被禁用的应用重新启用，相当于改变应用状态
func Enable(packageName string) {
	utils.Shell("pm enable " + packageName)
}

// IsEnabled 判断应用是否被启用
func IsEnabled(packageName string) bool {
	// 是否能正常启动，能启动则返回true
	return Launch(packageName, 0) == true
}

// GetBrowserPackage 获取系统默认浏览器包名
func GetBrowserPackage() string {
	str := utils.Shell("pm resolve-activity -a android.intent.action.VIEW -d http://example.com")
	re := regexp.MustCompile(`packageName=([a-zA-Z0-9_]+\.[^\s]+)`)
	match := re.FindStringSubmatch(str)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

// OpenUrl 用浏览器打开网站url。
func OpenUrl(url string) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	StartActivity(IntentOptions{
		Action: "VIEW",
		Data:   url,
	})
}

// StartActivity 根据选项构造一个Intent，并启动该Activity。
func StartActivity(options IntentOptions) {
	utils.Shell(buildIntentCommand(options, "start"))
}

// SendBroadcast 根据选项构造一个Intent，并发送该广播。
func SendBroadcast(options IntentOptions) {
	utils.Shell(buildIntentCommand(options, "broadcast"))
}

// StartService 根据选项构造一个Intent，并启动该服务。
func StartService(options IntentOptions) {
	utils.Shell(buildIntentCommand(options, "startservice"))
}

// BuildIntentCommand 根据选项构造一个Intent，并返回构造的命令字符串。
// 将高级的 Intent 配置转换为低级的 shell 命令，为上层函数提供基础支持。
// options：IntentOptions类型，包含Intent的选项
// commandType：字符串类型，表示Intent的类型，可选值为 "start"、"broadcast"、"startservice"
// 返回值：
// 构造的命令字符串
func buildIntentCommand(options IntentOptions, commandType string) string {
	var commandBuilder strings.Builder

	commandBuilder.WriteString("am " + commandType)

	if options.Action != "" {
		if strings.HasPrefix(options.Action, "android.intent.action.") {
			commandBuilder.WriteString(" -a " + options.Action)
		} else {
			commandBuilder.WriteString(" -a android.intent.action." + options.Action)
		}
	}

	if options.Type != "" {
		commandBuilder.WriteString(" -t " + options.Type)
	}

	if options.Data != "" {
		commandBuilder.WriteString(" -d " + options.Data)
	}

	for _, category := range options.Category {
		commandBuilder.WriteString(" -c " + category)
	}

	if options.PackageName != "" {
		commandBuilder.WriteString(" -n " + options.PackageName)
		if options.ClassName != "" {
			commandBuilder.WriteString("/" + options.ClassName)
		}
	}

	for key, value := range options.Extras {
		commandBuilder.WriteString(" --es " + key + " \"" + value + "\"")
	}

	for _, flag := range options.Flags {
		commandBuilder.WriteString(" --ez " + flag)
	}

	return commandBuilder.String()
}

// getMimeType 获取文件的MIME类型,为其他应用提供文件查看、编辑等功能。
func getMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	case ".apk":
		return "application/vnd.android.package-archive"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".flac":
		return "audio/flac"
	case ".aac":
		return "audio/aac"
	case ".m4a":
		return "audio/mp4"
	default:
		// Try to detect mime type based on file extension
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			return "*/*"
		}
		return mimeType
	}
}
