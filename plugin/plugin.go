package plugin

// 语言插件系统

import (
	"fmt"
	"github.com/xiaocainiao633/Genie1.0--/utils"
	"os"
	"strconv"
	"strings"
)

// 插件管理器
type Plugin struct {
	obj     string 
	apkPath string // apk文件路径
}

// java类包装
type Obj struct {
	obj string
}

// java对象包装
type Class struct {
	obj string
}

type Bitmap struct {
	Left, Top, Right, Bottom int
}

type BitmapFromBase64 struct {
	Base64 string
}

type BitmapFromPath struct {
	Path string
}

type AssetManager struct {
	ApkPath string
}

// 加载apk插件
func LoadApk(apkPath string) *Plugin {
	_, err := os.Stat(apkPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "File "+apkPath+" does not exist.")
		os.Exit(1)
	}
	str := utils.CallJavaMethod("plugin", "loadApk|"+apkPath)
	if str == "" {
		return nil
	}
	return &Plugin{obj: str, apkPath: apkPath}
}

// 创建类实例
func (p *Plugin) NewInstance(className string, values ...interface{}) *Class {
	list := "null"
	if len(values) > 0 {
		list = ""
		for _, v := range values {
			switch v.(type) {
			case string:
				list = list + "string@@" + v.(string) + "@@"
			case int:
				list = list + "int@@" + i2s(v.(int)) + "@@"
			case bool:
				list = list + "boolean@@" + b2s(v.(bool)) + "@@"
			case float32:
				list = list + "float@@" + fmt.Sprintf("%f", v.(float32)) + "@@"
			case float64:
				list = list + "double@@" + strconv.FormatFloat(v.(float64), 'f', -1, 64) + "@@"
			case int64:
				list = list + "long@@" + strconv.FormatInt(v.(int64), 10) + "@@"
			case AssetManager:
				list = list + "assetManager@@" + v.(AssetManager).ApkPath + "@@"
			case Bitmap:
				list = list + "bitmap@@" + i2s(v.(Bitmap).Left) + "," + i2s(v.(Bitmap).Top) + "," + i2s(v.(Bitmap).Right) + "," + i2s(v.(Bitmap).Bottom) + "@@"
			case BitmapFromBase64:
				list = list + "bitmapbase64@@" + v.(BitmapFromBase64).Base64 + "@@"
			case BitmapFromPath:
				list = list + "bitmappath@@" + v.(BitmapFromPath).Path + "@@"
			default:
				panic("Unsupported type: " + fmt.Sprintf("%T", v))
			}
		}
		list = list[:len(list)-2]
	}
	str := utils.CallJavaMethod("plugin", "newInstance|"+p.obj+"|"+className+"|"+list)
	if str == "" {
		return nil
	}
	return &Class{obj: str}
}

// 调用对象方法
func (c *Class) Call(methodName string, values ...interface{}) *Obj {
	list := "null"
	if len(values) > 0 {
		list = ""
		for _, v := range values {
			switch v.(type) {
			case string:
				s := v.(string)
				if s == "" {
					s = "null"
				}
				list = list + "string@@" + s + "@@"
			case int:
				list = list + "int@@" + i2s(v.(int)) + "@@"
			case bool:
				list = list + "boolean@@" + b2s(v.(bool)) + "@@"
			case float32:
				list = list + "float@@" + fmt.Sprintf("%f", v.(float32)) + "@@"
			case float64:
				list = list + "double@@" + strconv.FormatFloat(v.(float64), 'f', -1, 64) + "@@"
			case int64:
				list = list + "long@@" + strconv.FormatInt(v.(int64), 10) + "@@"
			case AssetManager:
				list = list + "assetManager@@" + v.(AssetManager).ApkPath + "@@"
			case Bitmap:
				list = list + "bitmap@@" + i2s(v.(Bitmap).Left) + "," + i2s(v.(Bitmap).Top) + "," + i2s(v.(Bitmap).Right) + "," + i2s(v.(Bitmap).Bottom) + "@@"
			case BitmapFromBase64:
				list = list + "bitmapbase64@@" + v.(BitmapFromBase64).Base64 + "@@"
			case BitmapFromPath:
				list = list + "bitmappath@@" + v.(BitmapFromPath).Path + "@@"
			default:
				panic("Unsupported type: " + fmt.Sprintf("%T", v))
			}
		}
		list = list[:len(list)-2]
	}
	return &Obj{obj: utils.CallJavaMethod("plugin", "call|"+c.obj+"|"+methodName+"|"+list)}
}

/*func (o *Obj) Release() {
	utils.CallJavaMethod("plugin", "releaseObj|"+o.obj)
	o.obj = ""
}*/

// 对象的一系列类型转换方法
func (o *Obj) ToString() string {
	return o.obj
}

func (o *Obj) ToInt() int {
	return s2i(o.obj)
}

func (o *Obj) ToBool() bool {
	return o.obj == "true"
}

func (o *Obj) ToInt64() int64 {
	if len(o.obj) > 1 {
		l, _ := strconv.ParseInt(o.obj, 10, 64)
		return l
	}
	return 0
}

func (o *Obj) ToFloat32() float32 {
	if len(o.obj) > 1 {
		f, _ := strconv.ParseFloat(o.obj, 32)
		return float32(f)
	}
	return 0
}

func s2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func i2s(i int) string {
	return strconv.Itoa(i)
}

func b2s(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
