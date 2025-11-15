package media

import (
	"github.com/xiaocainiao633/Genie1.0--/utils"
	"regexp"
	"strings"
)

// ScanFile 扫描路径path的媒体文件，将它加入媒体库中
func ScanFile(path string) {
	utils.Shell("am broadcast -a android.intent.action.MEDIA_SCANNER_SCAN_FILE -d \"file://" + path + "\"")
	mediaPath := strings.Replace(path, "/sdcard", "/storage/emulated/0", 1)
	result := utils.Shell(`content query --uri content://media/external/images/media`)
	lines := strings.Split(result, "\n")
	var mediaID string
	for _, line := range lines {
		if strings.Contains(line, "_data="+mediaPath) {
			re := regexp.MustCompile(`\b_id=([0-9]+)\b`)
			match := re.FindStringSubmatch(line)
			if len(match) == 2 {
				mediaID = match[1]
			}
			break
		}
	}
	if mediaID != "" {
		utils.Shell("content update --uri content://media/external/images/media/" + mediaID + " --bind is_pending:i:0")
	}
}

// PlayMP3 播放指定路径的 MP3 文件（使用系统播放器）
func PlayMP3(path string) {
	fileURI := "file://" + path
	utils.Shell(`am start -a android.intent.action.VIEW -d "` + fileURI + `" -t audio/*`)
}

// SendSMS 向指定手机号发送短信（打开短信界面，需用户确认）
func SendSMS(number, message string) {
	utils.Shell(`am start -a android.intent.action.SENDTO -d "smsto:` + number + `" --es sms_body "` + message + `"`)
}
