package uiacc

// 提供基于辅助功能服务的控件定位、交互操作等功能。无需开启APP的无障碍服务，通过程序化方式访问和操作界面元素
// 通过调用 Java 层的 JavaScript 引擎来执行传入的脚本，并返回执行结果
// 工作流程：
// 1. 创建一个 Uiacc 对象，用于构建选择器
// 2. 使用 Uiacc 对象的方法构建选择器
// 3. 使用 Uiacc 对象的方法查找控件
// 4. 使用 Uiacc 对象的方法点击控件
// 5. 使用 Uiacc 对象的方法长按控件
// 6. 使用 Uiacc 对象的方法选择控件
// 7. 使用 Uiacc 对象的方法清除选择控件
// 8. 使用 Uiacc 对象的方法设置选择控件


// 查找并点击文本为"确定"的按钮
// uiacc.New().Text("确定").FindOnce().Click()

// 等待元素出现并获取
// obj := uiacc.New().Id("button1").WaitFor(5000)

// 设置输入框文本
// uiacc.New().Editable(true).FindOnce().SetText("Hello")

import (
	_ "embed"
	"encoding/base64"
	"github.com/xiaocainiao633/Genie1.0--/motion"
	"github.com/xiaocainiao633/Genie1.0--/rhino"
	"github.com/xiaocainiao633/Genie1.0--/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Uiacc struct {
	selector string
}

type UiObject struct {
	index  int
	objStr string
}

type Rect struct {
	Left    int
	Right   int
	Top     int
	Bottom  int
	CenterX int
	CenterY int
	Width   int
	Height  int
}

//go:embed uiacc.js
var _uiacc_js string

var mutex sync.Mutex
var state bool
var index int

func init() {
	state = true
	utils.Shell("am broadcast -a com.autogo --es message refreshNode")
	utils.Shell("pkill -f assistdir")
	rhino.Eval("_node", _uiacc_js)
	node := &Uiacc{}
	node.PackageNameContains(".").WaitFor(500)
}

// New 创建一个新的 Accessibility 对象
func New() *Uiacc {
	mutex.Lock()
	defer mutex.Unlock()
	if !state {
		utils.Shell("am broadcast -a com.autogo --es message refreshNode")
		utils.Shell("pkill -f assistdir")
		rhino.Eval("_node", "init()")
		state = true
	}
	node := &Uiacc{}
	return node
}

// Text 设置选择器的 text 属性
func (a *Uiacc) Text(value string) *Uiacc {
	return &Uiacc{selector: a.selector + "text@@" + value + "&&"}
}

// TextContains 设置选择器的 textContains 属性，用于匹配包含指定文本的控件
func (a *Uiacc) TextContains(value string) *Uiacc {
	return &Uiacc{a.selector + "textContains@@" + value + "&&"}
}

// TextStartsWith 设置选择器的 textStartsWith 属性，用于匹配以指定文本开头的控件
func (a *Uiacc) TextStartsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "textStartsWith@@" + value + "&&"}
}

// TextEndsWith 设置选择器的 textEndsWith 属性，用于匹配以指定文本结尾的控件
func (a *Uiacc) TextEndsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "textEndsWith@@" + value + "&&"}
}

// TextMatches 设置选择器的 textMatches 属性，用于匹配符合指定正则表达式的控件
func (a *Uiacc) TextMatches(value string) *Uiacc {
	return &Uiacc{a.selector + "textMatches@@" + value + "&&"}
}

// Desc 设置选择器的 desc 属性，用于匹配描述等于指定文本的控件
func (a *Uiacc) Desc(value string) *Uiacc {
	return &Uiacc{a.selector + "desc@@" + value + "&&"}
}

// DescContains 设置选择器的 descContains 属性，用于匹配描述包含指定文本的控件
func (a *Uiacc) DescContains(value string) *Uiacc {
	return &Uiacc{a.selector + "descContains@@" + value + "&&"}
}

// DescStartsWith 设置选择器的 descStartsWith 属性，用于匹配描述以指定文本开头的控件
func (a *Uiacc) DescStartsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "descStartsWith@@" + value + "&&"}
}

// DescEndsWith 设置选择器的 descEndsWith 属性，用于匹配描述以指定文本结尾的控件
func (a *Uiacc) DescEndsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "descEndsWith@@" + value + "&&"}
}

// DescMatches 设置选择器的 descMatches 属性，用于匹配描述符合指定正则表达式的控件
func (a *Uiacc) DescMatches(value string) *Uiacc {
	return &Uiacc{a.selector + "descMatches@@" + value + "&&"}
}

// Id 设置选择器的 id 属性，用于匹配ID等于指定值的控件
func (a *Uiacc) Id(value string) *Uiacc {
	return &Uiacc{a.selector + "id@@" + value + "&&"}
}

// IdContains 设置选择器的 idContains 属性，用于匹配ID包含指定值的控件
func (a *Uiacc) IdContains(value string) *Uiacc {
	return &Uiacc{a.selector + "idContains@@" + value + "&&"}
}

// IdStartsWith 设置选择器的 idStartsWith 属性，用于匹配ID以指定值开头的控件
func (a *Uiacc) IdStartsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "idStartsWith@@" + value + "&&"}
}

// IdEndsWith 设置选择器的 idEndsWith 属性，用于匹配ID以指定值结尾的控件
func (a *Uiacc) IdEndsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "idEndsWith@@" + value + "&&"}
}

// IdMatches 设置选择器的 idMatches 属性，用于匹配ID符合指定正则表达式的控件
func (a *Uiacc) IdMatches(value string) *Uiacc {
	return &Uiacc{a.selector + "idMatches@@" + value + "&&"}
}

// ClassName 设置选择器的 className 属性，用于匹配类名等于指定值的控件
func (a *Uiacc) ClassName(value string) *Uiacc {
	return &Uiacc{a.selector + "className@@" + value + "&&"}
}

// ClassNameContains 设置选择器的 classNameContains 属性，用于匹配类名包含指定值的控件
func (a *Uiacc) ClassNameContains(value string) *Uiacc {
	return &Uiacc{a.selector + "classNameContains@@" + value + "&&"}
}

// ClassNameStartsWith 设置选择器的 classNameStartsWith 属性，用于匹配类名以指定值开头的控件
func (a *Uiacc) ClassNameStartsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "classNameStartsWith@@" + value + "&&"}
}

// ClassNameEndsWith 设置选择器的 classNameEndsWith 属性，用于匹配类名以指定值结尾的控件
func (a *Uiacc) ClassNameEndsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "classNameEndsWith@@" + value + "&&"}
}

// ClassNameMatches 设置选择器的 classNameMatches 属性，用于匹配类名符合指定正则表达式的控件
func (a *Uiacc) ClassNameMatches(value string) *Uiacc {
	return &Uiacc{a.selector + "classNameMatches@@" + value + "&&"}
}

// PackageName 设置选择器的 packageName 属性，用于匹配包名等于指定值的控件
func (a *Uiacc) PackageName(value string) *Uiacc {
	return &Uiacc{a.selector + "packageName@@" + value + "&&"}
}

// PackageNameContains 设置选择器的 packageNameContains 属性，用于匹配包名包含指定值的控件
func (a *Uiacc) PackageNameContains(value string) *Uiacc {
	return &Uiacc{a.selector + "packageNameContains@@" + value + "&&"}
}

// PackageNameStartsWith 设置选择器的 packageNameStartsWith 属性，用于匹配包名以指定值开头的控件
func (a *Uiacc) PackageNameStartsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "packageNameStartsWith@@" + value + "&&"}
}

// PackageNameEndsWith 设置选择器的 packageNameEndsWith 属性，用于匹配包名以指定值结尾的控件
func (a *Uiacc) PackageNameEndsWith(value string) *Uiacc {
	return &Uiacc{a.selector + "packageNameEndsWith@@" + value + "&&"}
}

// PackageNameMatches 设置选择器的 packageNameMatches 属性，用于匹配包名符合指定正则表达式的控件
func (a *Uiacc) PackageNameMatches(value string) *Uiacc {
	return &Uiacc{a.selector + "packageNameMatches@@" + value + "&&"}
}

// Bounds 设置选择器的 bounds 属性，用于匹配控件在屏幕上的范围
func (a *Uiacc) Bounds(left, top, right, bottom int) *Uiacc {
	return &Uiacc{a.selector + "bounds@@" + i2s(left) + "," + i2s(top) + "," + i2s(right) + "," + i2s(bottom) + "&&"}
}

// BoundsInside 设置选择器的 boundsInside 属性，用于匹配控件在屏幕内的范围
func (a *Uiacc) BoundsInside(left, top, right, bottom int) *Uiacc {
	return &Uiacc{a.selector + "boundsInside@@" + i2s(left) + "," + i2s(top) + "," + i2s(right) + "," + i2s(bottom) + "&&"}
}

// BoundsContains 设置选择器的 boundsContains 属性，用于匹配控件包含在指定范围内
func (a *Uiacc) BoundsContains(left, top, right, bottom int) *Uiacc {
	return &Uiacc{a.selector + "boundsContains@@" + i2s(left) + "," + i2s(top) + "," + i2s(right) + "," + i2s(bottom) + "&&"}
}

// DrawingOrder 设置选择器的 drawingOrder 属性，用于匹配控件在父控件中的绘制顺序
func (a *Uiacc) DrawingOrder(value int) *Uiacc {
	return &Uiacc{a.selector + "drawingOrder@@" + i2s(value) + "&&"}
}

// Clickable 设置选择器的 clickable 属性，用于匹配控件是否可点击
func (a *Uiacc) Clickable(value bool) *Uiacc {
	return &Uiacc{a.selector + "clickAble@@" + b2s(value) + "&&"}
}

// LongClickable 设置选择器的 longClickable 属性，用于匹配控件是否可长按
func (a *Uiacc) LongClickable(value bool) *Uiacc {
	return &Uiacc{a.selector + "longClickAble@@" + b2s(value) + "&&"}
}

// Checkable 设置选择器的 checkable 属性，用于匹配控件是否可选中
func (a *Uiacc) Checkable(value bool) *Uiacc {
	return &Uiacc{a.selector + "checkAble@@" + b2s(value) + "&&"}
}

// Selected 设置选择器的 selected 属性，用于匹配控件是否被选中
func (a *Uiacc) Selected(value bool) *Uiacc {
	return &Uiacc{a.selector + "selected@@" + b2s(value) + "&&"}
}

// Enabled 设置选择器的 enabled 属性，用于匹配控件是否启用
func (a *Uiacc) Enabled(value bool) *Uiacc {
	return &Uiacc{a.selector + "enabled@@" + b2s(value) + "&&"}
}

// Scrollable 设置选择器的 scrollable 属性，用于匹配控件是否可滚动
func (a *Uiacc) Scrollable(value bool) *Uiacc {
	return &Uiacc{a.selector + "scrollAble@@" + b2s(value) + "&&"}
}

// Editable 设置选择器的 editable 属性，用于匹配控件是否可编辑
func (a *Uiacc) Editable(value bool) *Uiacc {
	return &Uiacc{a.selector + "editable@@" + b2s(value) + "&&"}
}

// MultiLine 设置选择器的 multiLine 属性，用于匹配控件是否多行
func (a *Uiacc) MultiLine(value bool) *Uiacc {
	return &Uiacc{a.selector + "multiLine@@" + b2s(value) + "&&"}
}

// Checked 设置选择器的 checked 属性，用于匹配控件是否被勾选
func (a *Uiacc) Checked(value bool) *Uiacc {
	return &Uiacc{a.selector + "checked@@" + b2s(value) + "&&"}
}

// Focusable 设置选择器的 focusable 属性，用于匹配控件是否可聚焦
func (a *Uiacc) Focusable(value bool) *Uiacc {
	return &Uiacc{a.selector + "focusable@@" + b2s(value) + "&&"}
}

// Dismissable 设置选择器的 dismissable 属性，用于匹配控件是否可解散
func (a *Uiacc) Dismissable(value bool) *Uiacc {
	return &Uiacc{a.selector + "dismissable@@" + b2s(value) + "&&"}
}

// Focused 设置选择器的 UiaccFocused 属性，用于匹配控件是否是辅助功能焦点
func (a *Uiacc) Focused(value bool) *Uiacc {
	return &Uiacc{a.selector + "focused@@" + b2s(value) + "&&"}
}

// ContextClickable 设置选择器的 contextClickable 属性，用于匹配控件是否是上下文点击
func (a *Uiacc) ContextClickable(value bool) *Uiacc {
	return &Uiacc{a.selector + "contextClickable@@" + b2s(value) + "&&"}
}

// Index 设置选择器的 index 属性，用于匹配控件在父控件中的索引
func (a *Uiacc) Index(value int) *Uiacc {
	return &Uiacc{a.selector + "indexInParent@@" + i2s(value) + "&&"}
}

// Click 点击屏幕上的文本
func (a *Uiacc) Click(text string) bool {
	obj := a.Text(text).FindOnce()
	if obj != nil {
		return obj.Click() || obj.GetParent().Click()
	} else {
		obj = a.Desc(text).FindOnce()
		if obj != nil {
			return obj.Click() || obj.GetParent().Click()
		}
	}
	return false
}

// WaitFor 等待控件出现并返回 UiObject 对象 超时单位为毫秒,写0代表无限等待,超时返回nil
func (a *Uiacc) WaitFor(timeout int) *UiObject {
	startTime := time.Now()
	for {
		obj := a.FindOnce()
		if obj != nil {
			return obj
		}
		if timeout > 0 && time.Since(startTime).Milliseconds() >= int64(timeout) {
			break
		}
		sleep(100)
	}
	return nil
}

// FindOnce 查找单个控件并返回 UiObject 对象
func (a *Uiacc) FindOnce() *UiObject {
	return getNode("findOnce('" + a.selector + "');")
}

// Find 查找所有符合条件的控件并返回 UiObject 对象数组
func (a *Uiacc) Find() []*UiObject {
	mutex.Lock()
	defer mutex.Unlock()
	index++
	if index > 999 {
		index = 0
	}
	str := rhino.Eval("_node", "find("+i2s(index)+",'"+a.selector+"');")
	arr := strings.Split(str, "\n")
	if len(arr) < 2 { //因为返回值末尾带一个\n所以最小是两个成员
		return nil
	}
	var uiObjectArray []*UiObject
	for i := 0; i < len(arr)-1; i++ {
		uiObjectArray = append(uiObjectArray, &UiObject{objStr: arr[i], index: index})
		index++
		if index > 999 {
			index = 0
		}
	}
	return uiObjectArray
}

// Close 关闭无障碍服务
func Close() {
	mutex.Lock()
	defer mutex.Unlock()
	if state {
		rhino.Eval("_node", "close()")
		state = false
	}
}

// Click 点击该控件，并返回是否点击成功
func (u *UiObject) Click() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_CLICK);"))
}

// ClickCenter 使用坐标点击该控件的中点，相当于click(uiObj.bounds().centerX(), uiObject.bounds().centerY())
func (u *UiObject) ClickCenter() bool {
	bounds := u.GetBounds()
	if bounds.CenterX > 0 && bounds.CenterY > 0 {
		motion.Click(bounds.CenterX, bounds.CenterY, 1)
		return true
	}
	return false
}

// ClickLongClick 长按该控件，并返回是否点击成功
func (u *UiObject) ClickLongClick() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_LONG_CLICK);"))
}

// Copy 对输入框文本的选中内容进行复制，并返回是否操作成功
func (u *UiObject) Copy() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_COPY);"))
}

// Cut 对输入框文本的选中内容进行剪切，并返回是否操作成功
func (u *UiObject) Cut() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_CUT);"))
}

// Paste 对输入框控件进行粘贴操作，把剪贴板内容粘贴到输入框中，并返回是否操作成功
func (u *UiObject) Paste() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_PASTE);"))
}

// ScrollForward 对控件执行向前滑动的操作，并返回是否操作成功
func (u *UiObject) ScrollForward() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_SCROLL_FORWARD);"))
}

// ScrollBackward 对控件执行向后滑动的操作，并返回是否操作成功
func (u *UiObject) ScrollBackward() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_SCROLL_BACKWARD);"))
}

// Collapse 对控件执行折叠操作，并返回是否操作成功
func (u *UiObject) Collapse() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.AccessibilityAction.ACTION_COLLAPSE.getId();"))
}

// Expand 对控件执行展开操作，并返回是否操作成功
func (u *UiObject) Expand() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.AccessibilityAction.ACTION_EXPAND.getId();"))
}

// Show 执行显示操作，并返回是否操作成功
func (u *UiObject) Show() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.AccessibilityAction.ACTION_SHOW_ON_SCREEN.getId();"))
}

// Select 对控件执行"选中"操作，并返回是否操作成功
func (u *UiObject) Select() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_SELECT);"))
}

// ClearSelect 清除控件的选中状态，并返回是否操作成功
func (u *UiObject) ClearSelect() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_CLEAR_SELECTION);"))
}

// SetSelection 对输入框控件设置选中的文字内容，并返回是否操作成功
func (u *UiObject) SetSelection(start, end int) bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].setTextSelection("+i2s(start)+", "+i2s(end)+");"))
}

// SetVisibleToUser 设置控件是否可见
func (u *UiObject) SetVisibleToUser(isVisible bool) bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].setVisibleToUser("+b2s(isVisible)+");"))
}

// SetText 设置输入框控件的文本内容，并返回是否设置成功
func (u *UiObject) SetText(str string) bool {
	if str != "" {
		str = base64.StdEncoding.EncodeToString([]byte(str))
	}
	return s2b(rhino.Eval("_node", "var decodedBytes = Base64.decode('"+str+"', Base64.DEFAULT);var javaString = new java.lang.String(decodedBytes, 'UTF-8');var decodedText = String(javaString);var arguments = new Bundle();arguments.putCharSequence(AccessibilityNodeInfo.ACTION_ARGUMENT_SET_TEXT_CHARSEQUENCE, decodedText);nodeCache["+i2s(u.index)+"].performAction(AccessibilityNodeInfo.ACTION_SET_TEXT, arguments);"))
}

// GetClickable 获取控件的 clickable 属性
func (u *UiObject) GetClickable() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isClickable()"))
}

// GetLongClickable 获取控件的 longClickable 属性
func (u *UiObject) GetLongClickable() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isLongClickable()"))
}

// GetCheckable 获取控件的 checkable 属性
func (u *UiObject) GetCheckable() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isCheckable()"))
}

// GetSelected 获取控件的 selected 属性
func (u *UiObject) GetSelected() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isSelected()"))
}

// GetEnabled 获取控件的 enabled 属性
func (u *UiObject) GetEnabled() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isEnabled()"))
}

// GetScrollable 获取控件的 scrollable 属性
func (u *UiObject) GetScrollable() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isScrollable()"))
}

// GetEditable 获取控件的 editable 属性
func (u *UiObject) GetEditable() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isEditable()"))
}

// GetMultiLine 获取控件的 multiLine 属性
func (u *UiObject) GetMultiLine() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isMultiLine()"))
}

// GetChecked 获取控件的 checked 属性
func (u *UiObject) GetChecked() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isChecked()"))
}

// GetFocused 获取控件的 focused 属性
func (u *UiObject) GetFocused() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isFocused()"))
}

// GetFocusable 获取控件的 focusable 属性
func (u *UiObject) GetFocusable() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isFocusable()"))
}

// GetDismissable 获取控件的 dismissable 属性
func (u *UiObject) GetDismissable() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isDismissable()"))
}

// GetContextClickable 获取控件的 contextClickable 属性
func (u *UiObject) GetContextClickable() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isContextClickable()"))
}

// GetAccessibilityFocused 获取控件的 AccessibilityFocused 属性
func (u *UiObject) GetAccessibilityFocused() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isAccessibilityFocused()"))
}

// GetVisibleToUser 获取控件的 VisibleToUser 属性
func (u *UiObject) GetVisibleToUser() bool {
	return s2b(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].isVisibleToUser()"))
}

// GetChildCount 获取控件的子控件数目
func (u *UiObject) GetChildCount() int {
	return s2i(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].getChildCount()"))
}

// GetDrawingOrder 获取控件在父控件中的绘制次序
func (u *UiObject) GetDrawingOrder() int {
	return s2i(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].getDrawingOrder()"))
}

// GetIndex 获取控件在父控件中的索引
func (u *UiObject) GetIndex() int {
	js := `
(function(){
    var node = nodeCache[` + i2s(u.index) + `];
    
    var parent = node.getParent();
    if (!parent) return -1;
    
    var count = parent.getChildCount();
    for (var i = 0; i < count; i++) {
        var child = parent.getChild(i);
        if (child && child.equals(node)) {
            child.recycle();
            return i;
        }
        if (child) child.recycle();
    }
    return -1;
})()
`
	return s2i(rhino.Eval("_node", js))
}

// GetBounds 获取控件在屏幕上的范围
func (u *UiObject) GetBounds() Rect {
	str := rhino.Eval("_node", "var rect = new Rect();nodeCache["+i2s(u.index)+"].getBoundsInScreen(rect);rect.left + ',' + rect.top + ',' + rect.right + ',' + rect.bottom")
	arr := strings.Split(str, ",")
	if len(arr) != 4 {
		return Rect{}
	}
	return Rect{Left: s2i(arr[0]), Top: s2i(arr[1]), Right: s2i(arr[2]), Bottom: s2i(arr[3]), Width: s2i(arr[2]) - s2i(arr[0]), Height: s2i(arr[3]) - s2i(arr[1]), CenterX: s2i(arr[0]) + ((s2i(arr[2]) - s2i(arr[0])) / 2), CenterY: s2i(arr[1]) + ((s2i(arr[3]) - s2i(arr[1])) / 2)}
}

// GetBoundsInParent 获取控件在父控件中的范围
func (u *UiObject) GetBoundsInParent() Rect {
	str := rhino.Eval("_node", "var rect = new Rect();nodeCache["+i2s(u.index)+"].getBoundsInParent(rect);rect.left + ',' + rect.top + ',' + rect.right + ',' + rect.bottom")
	arr := strings.Split(str, ",")
	if len(arr) != 4 {
		return Rect{}
	}
	return Rect{Left: s2i(arr[0]), Top: s2i(arr[1]), Right: s2i(arr[2]), Bottom: s2i(arr[3]), Width: s2i(arr[2]) - s2i(arr[0]), Height: s2i(arr[3]) - s2i(arr[1]), CenterX: s2i(arr[0]) + ((s2i(arr[2]) - s2i(arr[0])) / 2), CenterY: s2i(arr[1]) + ((s2i(arr[3]) - s2i(arr[1])) / 2)}
}

// GetId 获取控件的资源ID
func (u *UiObject) GetId() string {
	js := `
(function(){
    var node = nodeCache[` + i2s(u.index) + `];
    
    var viewId = node.getViewIdResourceName();
    if (viewId != null) {
        var index = viewId.indexOf(":id/");
        if (index != -1) {
            return viewId.substring(index + 4);
        }
    }
    return "";
})()
`
	return s2s(rhino.Eval("_node", js))
}

// GetText 获取控件的文本内容
func (u *UiObject) GetText() string {
	return s2s(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].getText()"))
}

// GetDesc 获取控件的描述内容
func (u *UiObject) GetDesc() string {
	return s2s(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].getContentDescription()"))
}

// GetPackageName 获取控件的包名
func (u *UiObject) GetPackageName() string {
	return s2s(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].getPackageName()"))
}

// GetClassName 获取控件的类名
func (u *UiObject) GetClassName() string {
	return s2s(rhino.Eval("_node", "nodeCache["+i2s(u.index)+"].getClassName()"))
}

// GetParent 获取控件的父控件
func (u *UiObject) GetParent() *UiObject {
	return getNode("nodeCache[" + i2s(u.index) + "].getParent();")
}

// GetChild 获取控件的指定索引的子控件
func (u *UiObject) GetChild(index int) *UiObject {
	return getNode("nodeCache[" + i2s(u.index) + "].getChild(" + i2s(index) + ");")
}

// GetChildren 获取控件的所有子控件
func (u *UiObject) GetChildren() []*UiObject {
	mutex.Lock()
	defer mutex.Unlock()
	index++
	if index > 999 {
		index = 0
	}
	str := rhino.Eval("_node", "getChildren(nodeCache["+i2s(u.index)+"],"+i2s(index)+");")
	arr := strings.Split(str, "\n")
	if len(arr) < 2 { //因为返回值末尾带一个\n所以最小是两个成员
		return nil
	}
	var uiObjectArray []*UiObject
	for i := 0; i < len(arr)-1; i++ {
		uiObjectArray = append(uiObjectArray, &UiObject{objStr: arr[i], index: index})
		index++
		if index > 999 {
			index = 0
		}
	}
	return uiObjectArray
}

func getNode(js string) *UiObject {
	mutex.Lock()
	defer mutex.Unlock()
	index++
	if index > 999 {
		index = 0
	}
	str := rhino.Eval("_node", "nodeCache["+i2s(index)+"]="+js)
	if !strings.HasPrefix(str, "android.view.accessibility.AccessibilityNodeInfo@") {
		return nil
	}
	return &UiObject{objStr: str, index: index}
}

func b2s(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func s2b(s string) bool {
	return s == "true"
}

func i2s(i int) string {
	return strconv.Itoa(i)
}

func s2i(s string) int {
	i, _ := strconv.Atoi(strings.TrimSpace(s))
	return i
}

func s2s(s string) string {
	if s == "null" {
		s = ""
	}
	return s
}

func sleep(i int) {
	time.Sleep(time.Duration(i) * time.Millisecond)
}
