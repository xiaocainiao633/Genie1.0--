## 目录架构解析
```
opencv/
├── cgo.go              # CGO 配置和链接设置
├── opencv.go           # 高级封装：图像查找功能
├── core.h/core.go      # 核心数据结构（Mat、Point等）
├── core.cpp            # C++ 实现层
├── imgproc.h/imgproc.go # 图像处理函数
├── imgproc.cpp         # C++ 实现层
├── core_string.go      # 类型字符串转换
├── imgproc_colorcodes.go # 颜色空间转换常量
├── imgproc_string.go   # 图像处理相关字符串转换
└── opencv2/           # OpenCV 头文件
```

## 工作流程示例:
```
1. 用户调用 FindImage()
   ↓
2. opencv.go: byte2mat() 处理模板
   ├─ images.ReadFromBytes() 解码图像
   ├─ NewMatFromBytes() 创建 Mat
   ├─ checkTransparent() 检测透明
   ├─ createMask() 生成遮罩（如需要）
   ├─ matGray() 转灰度（如需要）
   └─ matScale() 缩放模板
   ↓
3. images.CaptureScreen() 截取屏幕
   ↓
4. NewMatFromBytes() 创建屏幕 Mat
   ↓
5. imgproc.go: MatchTemplate() 执行匹配
   ↓
6. core.cpp: C++ 调用 cv::matchTemplate()
   ↓
7. core.go: MinMaxLoc() 查找最佳匹配
   ↓
8. 返回匹配坐标
```