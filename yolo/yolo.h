#ifndef YOLO_H
#define YOLO_H

#ifdef __cplusplus
extern "C"
{
#endif

  // 封装模型及检测状态
  typedef struct Yolo Yolo;

  // 创建实例
  Yolo *newYolo();

  // 加载模型
  const char *loadModelYolo(Yolo *obj, const char *version, const char *param_path, const char *bin_path, const char *labels, int num_threads);

  // 目标检测
  char *detectYolo(Yolo *obj, const unsigned char *bitmapData, int width, int height, float nms, float prob, int size);

  // 关闭模型
  void closeYolo(Yolo *obj);

#ifdef __cplusplus
}
#endif

#endif // YOLO_H

// yolo目标检测库的c语言接口头文件
// 用于在图像中快速识别和检测目标对象