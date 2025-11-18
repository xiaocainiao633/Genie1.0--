# RAG 对话系统架构文档

## 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                   用户自然语言输入                         │
│              "点击登录按钮" / "输入文本'Hello'"            │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│              DialogueSystem (对话系统)                   │
│  - 接收用户输入                                          │
│  - 处理特殊命令 (help, exit)                            │
│  - 调用 Agent 处理查询                                   │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│                   Agent (核心Agent)                      │
│  - 协调各个组件                                          │
│  - 管理测试执行流程                                      │
│  - 格式化结果输出                                        │
└─────────────────────────────────────────────────────────┘
        ↓                           ↓
┌──────────────────┐      ┌──────────────────┐
│  CodeGenerator   │      │ KnowledgeBase    │
│  (代码生成器)     │      │ (知识库)          │
└──────────────────┘      └──────────────────┘
        ↓                           ↓
┌──────────────────┐      ┌──────────────────┐
│  Intent Parser   │      │  RAG Retrieval    │
│  (意图解析)      │      │  (检索)           │
└──────────────────┘      └──────────────────┘
        ↓                           ↓
┌─────────────────────────────────────────────────────────┐
│              代码生成 (Template + Rules)                 │
│  - 根据意图选择代码模板                                  │
│  - 填充参数                                              │
│  - 生成完整Go代码                                        │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│              代码编译和执行                              │
│  - 保存到文件                                            │
│  - Go编译                                                │
│  - (可选)执行测试                                        │
└─────────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────────┐
│              测试结果返回                                │
│  - 成功/失败状态                                         │
│  - 生成的代码                                            │
│  - 执行输出                                              │
│  - 错误信息                                              │
└─────────────────────────────────────────────────────────┘
```

## 核心组件

### 1. KnowledgeBase (知识库)

**职责**: 存储和检索API文档

**实现**:
- 使用 SQLite 数据库存储
- 支持关键词搜索
- 提供上下文检索（RAG）

**数据结构**:
```go
type APIDoc struct {
    Module      string  // 模块名
    Function    string  // 函数名
    Description string  // 描述
    Signature   string  // 函数签名
    Parameters  string  // 参数说明
    Return      string  // 返回值
    Example     string  // 示例代码
    Keywords    string  // 关键词（用于检索）
}
```

**关键方法**:
- `Search(query, limit)`: 搜索相关API
- `GetContext(query)`: 获取RAG上下文
- `AddAPI(doc)`: 添加API文档

### 2. CodeGenerator (代码生成器)

**职责**: 根据用户意图生成测试代码

**工作流程**:
1. **意图解析** (`parseIntent`):
   - 识别操作类型（click, input, assert等）
   - 识别目标类型（button, text, image等）
   - 提取参数值（坐标、文本、图片路径等）
   - 确定使用的模块

2. **代码生成** (`generateCode`):
   - 根据意图选择代码模板
   - 填充参数
   - 生成完整可执行代码

**支持的意图类型**:
- `click`: 点击操作
- `input`: 输入操作
- `assert`: 验证操作
- `wait`: 等待操作
- `launch`: 启动应用
- `swipe`: 滑动操作

### 3. Agent (核心Agent)

**职责**: 协调整个流程

**工作流程**:
1. 接收用户查询
2. 调用 CodeGenerator 生成代码
3. 保存代码到文件
4. 编译代码
5. (可选) 执行测试
6. 返回结果

### 4. DialogueSystem (对话系统)

**职责**: 提供交互式对话界面

**功能**:
- 交互式输入/输出
- 特殊命令处理（help, exit）
- 结果格式化显示

## 数据流

### 查询处理流程

```
用户输入: "点击登录按钮"
    ↓
[Intent Parser]
  - Action: "click"
  - Target: "button"
  - Value: "登录"
  - Module: "uiacc"
    ↓
[RAG Retrieval]
  - 搜索关键词: "点击", "button", "uiacc"
  - 检索相关API文档
  - 返回上下文信息
    ↓
[Code Generator]
  - 选择模板: generateClickCode
  - 填充参数: Text("登录")
  - 生成代码:
    obj := uiacc.New().Text("登录").FindOnce()
    if obj != nil {
        obj.Click()
    }
    ↓
[Code Compiler]
  - 保存到文件
  - Go编译
  - 检查编译错误
    ↓
[Result Formatter]
  - 格式化输出
  - 返回测试结果
```

## 扩展点

### 1. 添加新的意图类型

在 `code_generator.go` 中：

```go
// 在 parseIntent 中添加新的意图识别
if strings.Contains(queryLower, "新操作") {
    intent.Action = "new_action"
}

// 添加新的代码生成函数
func (cg *CodeGenerator) generateNewActionCode(intent Intent) string {
    // 生成代码逻辑
}
```

### 2. 集成本地LLM

修改 `code_generator.go` 的 `generateCode` 方法：

```go
func (cg *CodeGenerator) generateCodeWithLLM(intent Intent, context string) string {
    // 调用本地LLM (如Ollama)
    prompt := fmt.Sprintf(`
        根据以下上下文和用户意图生成Go测试代码:
        
        上下文:
        %s
        
        用户意图:
        %+v
        
        请生成完整的Go测试代码:
    `, context, intent)
    
    // 调用LLM API
    code := callLLM(prompt)
    return code
}
```

### 3. 增强检索能力

可以添加：
- **向量嵌入**: 使用sentence-transformers生成向量
- **语义搜索**: 使用向量相似度搜索
- **重排序**: 对检索结果进行重排序

### 4. 添加代码执行

在 `agent.go` 中添加：

```go
func (a *Agent) ExecuteTest(binaryPath string) (string, error) {
    // 在Android设备上执行
    cmd := exec.Command("adb", "shell", binaryPath)
    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

## 性能优化

1. **知识库缓存**: 缓存常用API文档
2. **代码模板缓存**: 缓存生成的代码模板
3. **并行检索**: 并行搜索多个关键词
4. **增量编译**: 只编译修改的部分

## 安全性

1. **代码验证**: 检查生成的代码安全性
2. **沙箱执行**: 在隔离环境中执行测试
3. **权限控制**: 限制可执行的操作
4. **输入验证**: 验证用户输入的有效性

