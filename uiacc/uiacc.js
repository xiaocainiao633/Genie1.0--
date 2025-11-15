// 导入所需的Android类
importClass(android.os.HandlerThread);
importClass(android.os.SystemClock);
importClass(android.app.UiAutomation);
importClass(android.app.UiAutomationConnection);
importClass(android.accessibilityservice.AccessibilityServiceInfo);
importClass(android.view.accessibility.AccessibilityEvent);
importClass(android.view.accessibility.AccessibilityNodeInfo);
importClass(android.graphics.Rect);
importClass(android.os.Build);
importClass(java.util.HashMap);
importClass(java.util.ArrayList);
importClass(java.util.regex.Pattern);
importClass(android.os.Bundle);
importClass(android.util.Base64);

// 全局变量
var HANDLER_THREAD_NAME = "AutoGo-UiDumpHandlerThread-" + SystemClock.uptimeMillis();
var mHandlerThread = new HandlerThread(HANDLER_THREAD_NAME);
var mUiAutomation = null;
var selectorMap = new HashMap();
var mInternalListener = null;
var lock = new Object();
var rootNode = null;
var cachedAllNodes = new ArrayList();
var nodesCacheValid = false; // 标记缓存是否有效
var nodeCache = new Array(1000)

// 初始化方法
function init() {
    // 检查是否需要重新初始化
    var needNewThread = false;

    if (mHandlerThread == null) {
        needNewThread = true;
    } else if (!mHandlerThread.isAlive()) {
        needNewThread = true;
    }

    // 如果需要新线程，创建并启动
    if (needNewThread) {
        var newThreadName = "AutoGo-UiDumpHandlerThread-" + SystemClock.uptimeMillis();
        mHandlerThread = new HandlerThread(newThreadName);

        mHandlerThread.start();
        mHandlerThread.getLooper();
        mUiAutomation = new UiAutomation(mHandlerThread.getLooper(), new UiAutomationConnection());

        try {
            mUiAutomation.connect();
        } catch (e) {
            throw new Error("设备中有其他无障碍服务正在运行导致出现冲突");
        }
    }

    // 确保其他组件正确初始化
    if (selectorMap == null) {
        selectorMap = new HashMap();
    } else {
        selectorMap.clear();
    }

    // 配置UiAutomation服务信息（每次都重新配置以确保正确）
    if (mUiAutomation != null) {
        var info = mUiAutomation.getServiceInfo();
        info.eventTypes = AccessibilityEvent.TYPES_ALL_MASK;
        info.feedbackType = AccessibilityServiceInfo.FEEDBACK_ALL_MASK;
        info.notificationTimeout = 100;
        info.flags = 122;
        info.packageNames = null;
        mUiAutomation.setServiceInfo(info);

        // 设置事件监听器
        setupEventListener();
    }
}

// 设置事件监听器
function setupEventListener() {
    mInternalListener = new UiAutomation.OnAccessibilityEventListener({
        onAccessibilityEvent: function (event) {
            if (event.getEventType() == AccessibilityEvent.TYPE_WINDOW_CONTENT_CHANGED) {
                synchronized(lock, function () {
                    // 回收旧的rootNode（如果存在）
                    if (rootNode != null) {
                        try {
                            rootNode.recycle();
                        } catch (e) {
                            // 静默处理
                        }
                    }

                    // 清理缓存的所有节点
                    clearNodesCache();

                    // 获取新的rootNode
                    try {
                        rootNode = mUiAutomation.getRootInActiveWindow();
                    } catch (e) {
                        rootNode = null;
                    }
                });
            }
        }
    });

    mUiAutomation.setOnAccessibilityEventListener(mInternalListener);
}

// synchronized 的 JavaScript 实现
function synchronized(lockObj, func) {
    // 在JavaScript中，我们使用简单的函数调用来模拟synchronized
    // 实际的同步需要依赖JavaScript运行环境的线程模型
    return func();
}

// 清理节点缓存
function clearNodesCache() {
    cachedAllNodes.clear();
    nodesCacheValid = false;
}

// 查找单个匹配节点
function findOnce(selectorStr) {
    // 检查是否已经关闭
    if (mUiAutomation == null) {
        return null;
    }

    if (selectorStr) {
        selector(selectorStr);
    }

    var allNodes = getCachedOrFreshNodes();

    try {
        for (var i = 0; i < allNodes.size(); i++) {
            var node = allNodes.get(i);
            if (hasNode(node)) {
                return node;
            }
        }
    } catch (e) {
        // 静默处理错误
    } finally {
        selectorMap.clear();
    }
    return null;
}

// 查找所有匹配节点
function find(index, selectorStr) {
    // 检查是否已经关闭
    if (mUiAutomation == null) {
        return "";
    }

    if (selectorStr) {
        selector(selectorStr);
    }

    var allNodes = getCachedOrFreshNodes();
    var str = ""
    try {
        for (var i = 0; i < allNodes.size(); i++) {
            var node = allNodes.get(i);
            if (hasNode(node)) {
                nodeCache[index] = node
                index = index + 1
                if (index > 999) {
                    index = 0
                }
                str = str + node.toString() + "\n"
            }
        }
    } catch (e) {
        // 静默处理错误
    } finally {
        selectorMap.clear();
    }
    return str;
}

// 获取指定控件的所有子控件
function getChildren(parentNode, index) {
    // 检查是否已经关闭
    if (mUiAutomation == null) {
        return "";
    }

    var str = "";
    try {
        var childCount = parentNode.getChildCount();
        for (var i = 0; i < childCount; i++) {
            var childNode = parentNode.getChild(i);
            if (childNode != null) {
                nodeCache[index] = childNode;
                index = index + 1;
                if (index > 999) {
                    index = 0;
                }
                str = str + "Child[" + i + "]: " + childNode.toString() + "\n";
            }
        }
    } catch (e) {
        // 静默处理错误
    }
    return str;
}

// 获取缓存的节点列表，如果缓存无效则重新遍历
function getCachedOrFreshNodes() {
    return synchronized(lock, function () {
        // 强制检查：获取当前根节点，与缓存进行比较
        var currentRoot = null;
        try {
            currentRoot = mUiAutomation.getRootInActiveWindow();
        } catch (e) {
            // 静默处理错误
        }

        // 如果当前根节点与缓存的根节点不同，强制更新缓存
        var needUpdate = false;
        if (currentRoot != rootNode) {
            needUpdate = true;
        }

        // 如果缓存有效且根节点没有变化，直接返回缓存
        if (nodesCacheValid && !cachedAllNodes.isEmpty() && !needUpdate) {
            return new ArrayList(cachedAllNodes); // 返回副本，避免外部修改
        }

        // 缓存无效或需要更新，重新遍历
        updateNodesCache();
        return new ArrayList(cachedAllNodes);
    });
}

// 更新节点缓存
function updateNodesCache() {
    // 先清理旧缓存
    clearNodesCache();

    // 获取当前最新的根节点，不依赖缓存的rootNode
    var currentRoot = null;
    try {
        currentRoot = mUiAutomation.getRootInActiveWindow();
    } catch (e) {
        return;
    }

    if (currentRoot != null) {
        // 使用当前根节点遍历
        traverseNode(currentRoot, cachedAllNodes);
        nodesCacheValid = true;
    } else {
        nodesCacheValid = false;
    }
}

// 关闭连接并释放所有资源
function close() {
    try {
        // 1. 移除事件监听器
        if (mUiAutomation != null && mInternalListener != null) {
            mUiAutomation.setOnAccessibilityEventListener(null);
        }

        // 2. 断开UiAutomation连接
        if (mUiAutomation != null) {
            mUiAutomation.disconnect();
        }

        // 3. 回收根节点
        if (rootNode != null) {
            try {
                rootNode.recycle();
            } catch (e) {
                // 静默处理
            }
            rootNode = null;
        }

        // 4. 清理缓存的节点列表
        if (cachedAllNodes != null) {
            cachedAllNodes.clear();
        }

        // 5. 清理节点缓存数组
        if (nodeCache != null) {
            for (var i = 0; i < nodeCache.length; i++) {
                nodeCache[i] = null;
            }
        }

        // 6. 清理选择器映射
        if (selectorMap != null) {
            selectorMap.clear();
        }

        // 7. 停止HandlerThread
        if (mHandlerThread != null && mHandlerThread.isAlive()) {
            mHandlerThread.quitSafely();
        }

        // 8. 重置状态标记和对象引用
        nodesCacheValid = false;
        mInternalListener = null;
        mUiAutomation = null;
        mHandlerThread = null;  // 重置为null，确保下次init()能正确创建新线程

    } catch (e) {
        // 静默处理错误
    }
}

// 设置选择器
function selector(str) {
    var arr = str.split("&&");
    for (var i = 0; i < arr.length; i++) {
        var s = arr[i].trim(); // 去掉首尾空格
        if (s == "") {
            continue;
        }

        var arrs = s.split("@@");
        if (arrs.length == 2 && arrs[0] != "") {
            var key = arrs[0].trim();
            var value = arrs[1].trim();
            selectorMap.put(key, value);
        }
    }
}

// 遍历节点
function traverseNode(nodeInfo, allNodes) {
    if (nodeInfo == null) {
        return;
    }
    allNodes.add(nodeInfo);

    for (var i = 0; i < nodeInfo.getChildCount(); i++) {
        var childNode = nodeInfo.getChild(i);
        traverseNode(childNode, allNodes);
    }
}

// 检查节点是否匹配选择器
function hasNode(nodeInfo) {
    if (nodeInfo == null) return false;

    var rect;
    var pattern;
    var list = []; // 使用JavaScript数组
    var entrySet = selectorMap.entrySet();
    var iterator = entrySet.iterator();

    while (iterator.hasNext()) {
        var entry = iterator.next();
        var key = String(entry.getKey()); // 转换为JavaScript字符串
        var value = String(entry.getValue()); // 转换为JavaScript字符串

        switch (key) {
            case "text":
                var nodeText = safeCharSeqToString(nodeInfo.getText());
                list.push(nodeText == value);
                break;
            case "textContains":
                var nodeText = safeCharSeqToString(nodeInfo.getText());
                list.push(nodeText.indexOf(value) != -1);
                break;
            case "textStartsWith":
                var nodeText = safeCharSeqToString(nodeInfo.getText());
                list.push(nodeText.indexOf(value) == 0);
                break;
            case "textEndsWith":
                var nodeText = safeCharSeqToString(nodeInfo.getText());
                list.push(nodeText.length >= value.length &&
                    nodeText.substring(nodeText.length - value.length) == value);
                break;
            case "textMatches":
                pattern = Pattern.compile(value);
                list.push(pattern.matcher(safeCharSeqToString(nodeInfo.getText())).matches());
                break;
            case "desc":
                var nodeDesc = safeCharSeqToString(nodeInfo.getContentDescription());
                list.push(nodeDesc == value);
                break;
            case "descContains":
                var nodeDesc = safeCharSeqToString(nodeInfo.getContentDescription());
                list.push(nodeDesc.indexOf(value) != -1);
                break;
            case "descStartsWith":
                var nodeDesc = safeCharSeqToString(nodeInfo.getContentDescription());
                list.push(nodeDesc.indexOf(value) == 0);
                break;
            case "descEndsWith":
                var nodeDesc = safeCharSeqToString(nodeInfo.getContentDescription());
                list.push(nodeDesc.length >= value.length &&
                    nodeDesc.substring(nodeDesc.length - value.length) == value);
                break;
            case "descMatches":
                pattern = Pattern.compile(value);
                list.push(pattern.matcher(safeCharSeqToString(nodeInfo.getContentDescription())).matches());
                break;
            case "id":
                var nodeId = getId(nodeInfo);
                list.push(nodeId == value);
                break;
            case "idContains":
                var nodeId = getId(nodeInfo);
                list.push(nodeId.indexOf(value) != -1);
                break;
            case "idStartsWith":
                var nodeId = getId(nodeInfo);
                list.push(nodeId.indexOf(value) == 0);
                break;
            case "idEndsWith":
                var nodeId = getId(nodeInfo);
                list.push(nodeId.length >= value.length &&
                    nodeId.substring(nodeId.length - value.length) == value);
                break;
            case "idMatches":
                pattern = Pattern.compile(value);
                list.push(pattern.matcher(getId(nodeInfo)).matches());
                break;
            case "className":
                var nodeClassName = safeCharSeqToString(nodeInfo.getClassName());
                list.push(nodeClassName == value);
                break;
            case "classNameContains":
                var nodeClassName = safeCharSeqToString(nodeInfo.getClassName());
                list.push(nodeClassName.indexOf(value) != -1);
                break;
            case "classNameStartsWith":
                var nodeClassName = safeCharSeqToString(nodeInfo.getClassName());
                list.push(nodeClassName.indexOf(value) == 0);
                break;
            case "classNameEndsWith":
                var nodeClassName = safeCharSeqToString(nodeInfo.getClassName());
                list.push(nodeClassName.length >= value.length &&
                    nodeClassName.substring(nodeClassName.length - value.length) == value);
                break;
            case "classNameMatches":
                pattern = Pattern.compile(value);
                list.push(pattern.matcher(safeCharSeqToString(nodeInfo.getClassName())).matches());
                break;
            case "packageName":
                var nodePackageName = safeCharSeqToString(nodeInfo.getPackageName());
                list.push(nodePackageName == value);
                break;
            case "packageNameContains":
                var nodePackageName = safeCharSeqToString(nodeInfo.getPackageName());
                list.push(nodePackageName.indexOf(value) != -1);
                break;
            case "packageNameStartsWith":
                var nodePackageName = safeCharSeqToString(nodeInfo.getPackageName());
                list.push(nodePackageName.indexOf(value) == 0);
                break;
            case "packageNameEndsWith":
                var nodePackageName = safeCharSeqToString(nodeInfo.getPackageName());
                list.push(nodePackageName.length >= value.length &&
                    nodePackageName.substring(nodePackageName.length - value.length) == value);
                break;
            case "packageNameMatches":
                pattern = Pattern.compile(value);
                list.push(pattern.matcher(safeCharSeqToString(nodeInfo.getPackageName())).matches());
                break;
            case "bounds":
                rect = new Rect();
                nodeInfo.getBoundsInScreen(rect);
                var rectStr = rect.left + "," + rect.top + "," + rect.right + "," + rect.bottom;
                list.push(rectStr == value);
                break;
            case "boundsInside":
                rect = new Rect();
                nodeInfo.getBoundsInScreen(rect);
                var parts = value.split(",");
                if (parts.length == 4) {
                    var left = parseInt(parts[0]);
                    var top = parseInt(parts[1]);
                    var right = parseInt(parts[2]);
                    var bottom = parseInt(parts[3]);
                    var isInside = rect.left < rect.right && rect.top < rect.bottom &&
                        rect.left >= left && rect.top >= top &&
                        rect.right <= right && rect.bottom <= bottom;
                    list.push(isInside);
                } else {
                    list.push(false);
                }
                break;
            case "boundsContains":
                rect = new Rect();
                nodeInfo.getBoundsInScreen(rect);
                var containsParts = value.split(",");
                if (containsParts.length == 4) {
                    var containsLeft = parseInt(containsParts[0]);
                    var containsTop = parseInt(containsParts[1]);
                    var containsRight = parseInt(containsParts[2]);
                    var containsBottom = parseInt(containsParts[3]);
                    var contains = rect.left < rect.right && rect.top < rect.bottom &&
                        rect.left <= containsLeft && rect.top <= containsTop &&
                        rect.right >= containsRight && rect.bottom >= containsBottom;
                    list.push(contains);
                } else {
                    list.push(false);
                }
                break;
            case "drawingOrder":
                list.push(Build.VERSION.SDK_INT >= Build.VERSION_CODES.N &&
                    nodeInfo.getDrawingOrder() == s2i(value));
                break;
            case "clickAble":
                list.push(nodeInfo.isClickable() == Boolean.parseBoolean(value));
                break;
            case "longClickAble":
                list.push(nodeInfo.isLongClickable() == Boolean.parseBoolean(value));
                break;
            case "checkAble":
                list.push(nodeInfo.isCheckable() == Boolean.parseBoolean(value));
                break;
            case "selected":
                list.push(nodeInfo.isSelected() == Boolean.parseBoolean(value));
                break;
            case "enabled":
                list.push(nodeInfo.isEnabled() == Boolean.parseBoolean(value));
                break;
            case "scrollAble":
                list.push(nodeInfo.isScrollable() == Boolean.parseBoolean(value));
                break;
            case "editable":
                list.push(nodeInfo.isEditable() == Boolean.parseBoolean(value));
                break;
            case "multiLine":
                list.push(nodeInfo.isMultiLine() == Boolean.parseBoolean(value));
                break;
            case "checked":
                list.push(nodeInfo.isChecked() == Boolean.parseBoolean(value));
                break;
            case "focusable":
                list.push(nodeInfo.isFocusable() == Boolean.parseBoolean(value));
                break;
            case "dismissable":
                list.push(nodeInfo.isDismissable() == Boolean.parseBoolean(value));
                break;
            case "contextClickable":
                list.push(Build.VERSION.SDK_INT >= Build.VERSION_CODES.M &&
                    nodeInfo.isContextClickable() == Boolean.parseBoolean(value));
                break;
            case "focused":
                list.push(nodeInfo.isAccessibilityFocused() == Boolean.parseBoolean(value));
                break;
            case "indexInParent":
                if (nodeInfo.getParent() != null) {
                    var parentNode = nodeInfo.getParent();
                    var found = false;
                    var indexInParent = -1;
                    var childCount = parentNode.getChildCount();
                    for (var j = 0; j < childCount; j++) {
                        var childNode = parentNode.getChild(j);
                        if (childNode != null && childNode.equals(nodeInfo)) {
                            indexInParent = j;
                            found = true;
                            childNode.recycle();
                            break;
                        }
                        if (childNode != null) {
                            childNode.recycle();
                        }
                    }
                    list.push(found && indexInParent == s2i(value));
                } else {
                    list.push(false);
                }
                break;
            default:
                // 未知的选择器类型，静默忽略
                break;
        }
    }

    // 如果list为空但selectorMap不为空，说明没有匹配到任何条件
    if (list.length == 0 && selectorMap.size() > 0) {
        return false;
    }

    // 检查条件数量是否匹配
    if (selectorMap.size() != list.length) {
        return false;
    }

    for (var i = 0; i < list.length; i++) {
        if (!list[i]) {
            return false;
        }
    }

    return true;
}

// 工具方法：获取ID
function getId(nodeInfo) {
    var viewId = nodeInfo.getViewIdResourceName();
    if (viewId != null) {
        var index = viewId.indexOf(":id/");
        if (index != -1) {
            return viewId.substring(index + 4);
        }
    }
    return "";
}

// 工具方法：字符串转整数
function s2i(s) {
    if (s != null && s.length > 0) {
        try {
            return parseInt(s);
        } catch (e) {
            // 忽略异常
        }
    }
    return 0;
}

// 工具方法：安全地将CharSequence转换为JavaScript字符串
function safeCharSeqToString(charSequence) {
    return charSequence == null ? "" : String(charSequence);
}

// 执行初始化
init();