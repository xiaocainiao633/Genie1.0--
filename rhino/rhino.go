package rhino

// 提供javascript执行能力并返回其结果

import (
	"github.com/xiaocainiao633/Genie1.0--/utils"
)

// 通过调用 Java 层的 JavaScript 引擎来执行传入的脚本，并返回执行结果
func Eval(contextId, script string) string {
	if script == "" {
		return ""
	}
	if contextId == "" {
		contextId = "__TEMP__"
	}
	return utils.CallJavaMethod("js", contextId+"|"+script)
}
