# AutoGo 自动化测试示例

本目录包含 AutoGo 框架的自动化测试示例脚本。

## 文件说明

- `login_interface_test.go` - 登录界面自动化测试示例

## 使用方法

### 1. 准备测试环境

```bash
# 确保已安装 Go 1.24.2 或更高版本
go version

# 安装依赖
go mod tidy
```

### 2. 准备测试资源

创建以下目录结构：

```
project/
├── templates/              # 图像模板目录
│   ├── login_page.png      # 登录页面模板
│   ├── login_button.png     # 登录按钮模板
│   ├── error_dialog.png    # 错误弹窗模板
│   └── login_failed.png    # 登录失败提示模板
└── examples/
    └── login_interface_test.go
```

### 3. 修改配置

在 `login_interface_test.go` 中修改应用包名：

```go
// 将 "com.example.app" 替换为实际的应用包名
app.Launch("com.example.app", 0)
```

### 4. 运行测试

```bash
# 编译
go build -o login_test examples/login_interface_test.go

# 在 Android 设备上运行
adb push login_test /data/local/tmp/
adb shell chmod +x /data/local/tmp/login_test
adb shell /data/local/tmp/login_test
```

## 测试用例说明

### TestUIElements
- **目的**: 验证登录界面的所有必需UI元素是否存在
- **验证项**: 用户名输入框、密码输入框、登录按钮

### TestErrorHandling
- **目的**: 测试错误处理机制
- **操作**: 输入错误的用户名和密码，验证是否显示错误提示

### TestLoginInterface
- **目的**: 完整的登录流程测试
- **步骤**:
  1. 启动应用
  2. 验证登录页面
  3. 输入用户名和密码
  4. 点击登录按钮
  5. 验证登录结果

## 测试结果

测试完成后会生成详细的测试报告，包括：
- 每个测试用例的执行状态（通过/失败）
- 测试消息
- 执行耗时
- 总体统计信息

## 注意事项

1. **图像模板**: 确保模板图片与实际界面匹配
2. **应用包名**: 根据实际应用修改包名
3. **等待时间**: 根据应用响应速度调整 `utils.Sleep()` 的时间
4. **相似度阈值**: 根据实际情况调整 `opencv.FindImage()` 的相似度参数

## 扩展测试

可以参考 `login_interface_test.go` 的结构，创建其他界面的测试脚本：
- 购物流程测试
- 表单填写测试
- 导航测试
- 数据验证测试

## 更多信息

详细 API 配合机制请参考: `../API_COOPERATION_GUIDE.md`

