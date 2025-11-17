# Genie1.0--
### 基于GO语言实现的安卓自动化测试项目，此项目即将完成

包含大量的api，使用RAG构建知识库，有内置agent可以根据输入的命令自动生成调用相关api的脚本进行测试，并输出测试结果，通过opencv、yolo等智能方式判断是否存在异常
支持自定义api进行测试，会自动将新的api加入知识库，具体使用方式见agent目录下的markdown文件或者下方的描述，核心文件是：opencv、agent、examples等

```
Genie1.0--/
├── app          # 应用相关API文件目录
├── console      # 控制台功能API文件目
├── device       # 设备相关API文件目录
├── files        # 文件相关API文件目录
├── https        # 网络请求相关API文件目录
├── images       # 图片资源API文件目录
├── ime          # 输入法相关API文件目录
├── imgui        # ImGui界面相关API文件目录
├── libs         # 依赖库文件目录
├── media        # 媒体资源相关API文件目录
├── motion       # 动效相关API文件目录
├── plugin       # 插件相关API文件目录
├── ppocr        # PPOCR文字识别相关API文件目录
├── rhino        # Rhino服务相关API文件目录
├── storages     # 存储功能相关API文件目录
├── system       # 系统相关API文件目录
├── uiacc        # UI交互相关API文件目录
├── utils        # 工具类相关API文件目录
├── workspace    # 工作区相关API文件目录
└── yolo         # YOLO目标检测相关API文件目录
```
