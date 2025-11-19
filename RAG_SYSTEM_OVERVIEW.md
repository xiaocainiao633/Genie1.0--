# AutoGo RAG 系统概览

## 概述

AutoGo RAG 系统 = “自然语言 → RAG 检索 → 本地 LLM → AutoGo API → Android 执行”。它可以自动理解测试需求、生成 Go 测试脚本、构建并（可选）推送到 Android 设备执行，最终输出详细的测试报告。

## 关键能力

- **本地 LLM（Ollama）**：默认集成 `llama3.2:latest`
- **向量知识库**：SQLite + Embedding（自动生成向量，支持语义检索）
- **多轮对话记忆**：记住最近 N 轮上下文，连续生成脚本
- **全链路自动化**：脚本生成 → `go build` →（可选）ADB push 执行
- **测试报告**：JSON + Markdown，包含编译/运行日志和异常明细

## 工作流程

```
用户输入 + 历史上下文
          ↓
DialogueSystem（多轮记忆）
          ↓
Agent
  ├─ CodeGenerator：规则模板 + LLM 生成
  ├─ KnowledgeBase：SQLite + 向量检索
  ├─ AndroidExecutor：ADB push / shell 执行
  └─ ReportWriter：JSON / Markdown 报告
          ↓
Go 脚本 → 可执行文件 → Android 设备运行
```

## 自动化测试原理

1. **自然语言解析**：识别操作类型（点击/输入/验证/等待/滑动…）及目标、参数
2. **RAG 检索**：从知识库中检索相关 API 文档 & 示例（关键词 + 向量双通道）
3. **代码生成**：
   - 模板引擎覆盖 80% 常见操作
   - LLM 将检索结果 + 上下文作为 Prompt，生成完整脚本
4. **Go 编译**：脚本写入 `workspace` 并 `go build`
5. **自动执行（可选）**：ADB push 到 `/data/local/tmp`，chmod + shell 执行
6. **报告输出**：收集编译输出、设备日志、错误信息，生成 JSON/Markdown 报告

> 综上，本项目是通过“自动生成 AutoGo API 脚本”的方式进行自动化测试，而不是依赖传统的手写脚本。

## 页面异常判定机制

| 技术 | API | 用途 | 异常示例 |
|------|-----|------|---------|
| 模板匹配 | `opencv.FindImage` | 对比布局/按钮/颜色 | 元素缺失、颜色异常、布局错位 |
| 目标检测 | `yolo.Detect` | 识别动态异常元素 | 异常弹窗、广告遮挡 |
| OCR 识别 | `ppocr.Ocr` | 捕捉提示语 | “失败”“错误”“网络异常” |
| 颜色检测 | `images.FindColor` | 监控状态条颜色 | Loading 条长时间红色 |
| UI 控件 | `uiacc.New().Text(...).FindOnce()` | 获取控件状态 | 控件不可见/不可点击 |

若任一检测项失败，系统会在测试报告中记录“触发异常的步骤 + 证据”，并将该轮测试标记为失败。

## 目录结构

```
agent/
├── agent.go             # 调度/编译/执行/报告
├── code_generator.go    # 规则模板 + LLM
├── knowledge_base.go    # SQLite + Embedding
├── ollama_client.go     # Ollama HTTP 客户端
├── memory.go            # 多轮对话记忆
├── executor.go          # ADB 推送执行
├── report.go            # 报告结构
├── dialogue.go          # 对话系统
└── config.go            # 配置项
```

根目录 `main.go` 提供 CLI：

```bash
go run main.go \
  -kb ./knowledge_base.db \
  -workspace ./workspace \
  -use-llm=true \
  -auto-exec=false \
  -query "点击登录按钮"
```

## 简单理解几个问题：详细请查看下方的相关文档

### 1. 此项目是如何进行自动化测试的？

- RAG + LLM 根据自然语言生成 Go 测试脚本
- 脚本内部调用 AutoGo API（motion/uiacc/opencv/ppocr/yolo/images…）完成操作与断言
- 编译后可直接推送到 Android 设备执行
- 收集日志 → 生成报告 → 返回给用户

### 2. 怎样判断页面是否异常？

系统组合多种信号：模板匹配（OpenCV）、目标检测（YOLO）、文字识别（PPOCR）、颜色检测（Images）、UI 属性（UIACC）。只要任一环节检测到异常，如“未匹配到按钮模板”“OCR 识别出‘错误’字样”，就会在报告中记录并标记失败。

## 相关文档

- [ agent 目录概述](./agent/README.md)
- [ agent 使用方式详解](./agent/USAGE_GUIDE.md)
- [ agent 整体架构分析](./agent/ARCHITECTURE.md)
- [ 项目快速开始 ](./agent/QUICK_START.md)

---
