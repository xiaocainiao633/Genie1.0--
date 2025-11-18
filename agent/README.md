# AutoGo RAG 对话系统

基于 RAG（Retrieval-Augmented Generation）与本地 LLM 的自动化测试脚本生成系统。

## 功能特性

- 🧠 **本地 LLM**：默认对接 Ollama `llama3.2:latest`，支持离线推理
- 🗂 **RAG 知识库**：SQLite + 向量检索（自动生成 Embedding）
- 🔁 **多轮记忆**：对话系统自动维护最近 N 轮上下文
- ⚙️ **自动脚本生成**：根据自然语言自动组合 AutoGo API（motion/uiacc/opencv/ppocr/yolo…）
- 📱 **ADB 自动执行**：一键 push 到 Android 设备执行
- 📝 **测试报告**：输出 JSON + Markdown 报告，包含编译/运行日志

## 系统架构

```
┌──────────────────────────────────────────────────────┐
│             用户自然语言 / RAG 记忆上下文              │
└──────────────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────┐
│      DialogueSystem（多轮对话 + 上下文记忆）           │
└──────────────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────┐
│              Agent（调度编排核心）                     │
│  ├─ CodeGenerator（规则模板 + LLM 生成）              │
│  ├─ KnowledgeBase（SQLite + 向量检索）                 │
│  ├─ AndroidExecutor（ADB push & shell）               │
│  └─ ReportWriter（JSON / Markdown 报告）              │
└──────────────────────────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────┐
│              测试脚本 → 可执行文件 → 设备运行            │
└──────────────────────────────────────────────────────┘
```

## 自动化测试原理

1. **自然语言解析**：识别操作意图（点击/输入/验证/滑动…）
2. **RAG 检索**：根据意图到知识库检索相关 API 文档与示例
3. **代码生成**：
   - 模板引擎：命中常见意图时直接拼接测试脚本
   - LLM：复杂场景会把检索到的 API + 对话上下文一起交给本地 LLM
4. **脚本编译**：生成的 Go 代码被编译成二进制
5. **自动执行（可选）**：通过 `adb push` 上传至 Android 设备执行
6. **报告输出**：记录编译日志、设备运行日志、错误堆栈，生成 JSON/Markdown 报告

> 整个过程无需手写脚本，系统会自动组合 AutoGo SDK（motion/uiacc/opencv/ppocr/yolo/images 等）完成操作与断言。

## 页面异常判定机制

| 技术 | 作用 | 异常判定示例 |
|------|------|-------------|
| **OpenCV** (`opencv.FindImage`) | 模板匹配 | 判断按钮/模块是否缺失，检测布局错位 |
| **YOLO** (`yolo.Detect`) | 目标检测 | 识别误弹窗、广告遮挡、异常浮层 |
| **OCR** (`ppocr.Ocr`) | 文本识别 | 捕获“失败”“错误”“断网”等提示词 |
| **颜色检测** (`images.FindColor`) | 状态指示 | 校验提示条颜色、Loading 状态 |
| **UI 无障碍** (`uiacc`) | 属性校验 | 获取控件状态（enabled/visible/focused） |

系统会按照“模板匹配 → YOLO → OCR → UI 属性”多级校验，若任一环节检测到异常，则写入测试报告并标记失败原因。

## 快速开始

```bash
# 1. 初始化知识库
go run main.go -init -kb ./knowledge_base.db

# 2. 启动对话系统
go run main.go -kb ./knowledge_base.db -workspace ./workspace

# 3. 单次查询
go run main.go -kb ./knowledge_base.db -workspace ./workspace \
  -query "启动应用com.example.app并点击登录按钮"
```

常用参数：

| 参数 | 说明 |
|------|------|
| `-use-llm=true` | 启用 Ollama 生成（默认开启） |
| `-auto-exec=true` | 自动 push 到设备执行 |
| `-adb` | 自定义 adb 路径 |
| `-remote` | 设备端执行目录（默认 `/data/local/tmp`） |

## 示例意图

- `点击登录按钮`
- `在输入框中输入'Hello World'`
- `验证页面是否存在'主页'文字'`
- `启动应用com.example.app`
- `从(100,200)滑动到(300,800)`
- `检查是否弹出“支付失败”`

## 知识库扩展

```go
doc := agent.APIDoc{
    Module:   "motion",
    Function: "NewAction",
    // ...
}
kb.AddAPI(doc)
```

系统会自动为新文档生成 Embedding，以便向量检索。

## 项目结构

```
agent/
├── agent.go             # 调度核心（编译/执行/报告）
├── code_generator.go    # 意图解析 + 模板/LLM生成
├── dialogue.go          # 多轮对话系统
├── knowledge_base.go    # SQLite + 向量检索
├── memory.go            # 对话记忆
├── executor.go          # ADB 执行器
├── report.go            # 测试报告
├── ollama_client.go     # Ollama 接入
└── config.go            # 运行配置
```

## 常见问题

### 1. 系统是怎样进行自动化测试的？

RAG 解析需求 → 自动生成 Go 脚本 → 编译为二进制 →（可选）ADB 上传到 Android 设备执行 → 汇总日志/截图 → 输出报告。整个脚本都是由 AutoGo 的 API 组合而成。

### 2. 如何判断页面是否异常？

通过“模板匹配 + 目标检测 + OCR + UI 属性”多层校验。典型流程：

1. 用 `opencv.FindImage` 校验 UI 布局/按钮位置
2. 用 `yolo.Detect` 检测异常弹窗、广告、遮挡
3. 用 `ppocr.Ocr` 识别错误提示/失败文本
4. 用 `uiacc` 读取控件状态（可点击/可见/可输入）
5. 任何一步异常都会被写入测试报告

## 运行前置条件

1. **Go 1.24.2+**
2. **Ollama**（可选，默认使用 `llama3.2:latest`）
3. **ADB**（自动执行需要连接 Android 设备或模拟器）
4. **AutoGo SDK**（当前仓库）

## 相关文档

- [ARCHITECTURE.md](./ARCHITECTURE.md) - 详细架构
- [USAGE_GUIDE.md](./USAGE_GUIDE.md) - 使用指南 & 问题解答
- [QUICK_START.md](./QUICK_START.md) - 5 分钟上手
- [RAG_SYSTEM_OVERVIEW.md](../RAG_SYSTEM_OVERVIEW.md) - 总体概览

