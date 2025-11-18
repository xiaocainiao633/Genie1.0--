# 快速开始指南

### 1. 初始化系统

```bash
# 进入项目目录
cd AutoGo

# 初始化知识库
go run main.go -init -kb ./knowledge_base.db
```

### 2. 启动对话系统

```bash
go run main.go -kb ./knowledge_base.db -workspace ./workspace
```

### 3. 开始使用

输入自然语言描述，系统会自动生成测试代码：

```
你: 点击登录按钮

[系统] 正在处理您的请求...
[Agent] 正在分析用户需求...
[Agent] 代码生成完成

生成的代码:
---
package main
...
---
```

## 常用命令

- `help` - 显示帮助信息
- `exit` / `quit` - 退出系统

## 示例查询

### 基础操作

- `点击登录按钮`
- `在坐标(100, 200)点击`
- `输入文本'Hello World'`
- `等待5秒`

### 验证操作

- `验证页面是否存在'主页'文字`
- `验证是否存在图片logo.png`

### 应用操作

- `启动应用com.example.app`
- `滑动从(100,200)到(300,400)`

## 下一步

- 阅读 [USAGE_GUIDE.md](./USAGE_GUIDE.md) 了解详细用法
- 阅读 [ARCHITECTURE.md](./ARCHITECTURE.md) 了解系统架构
- 查看 [README.md](./README.md) 了解完整功能
