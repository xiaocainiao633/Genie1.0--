#ifndef PPOCR_H
#define PPOCR_H

#ifdef __cplusplus
extern "C"
{
#endif

  typedef struct Ppocr Ppocr;

  Ppocr *newPpocr();

  // 加载 OCR 识别所需的模型文件
  const char *loadModelPpocr(Ppocr *Obj, const char *label, const char *dbParam, const char *dbBin, const char *recParam, const char *recBin, int thread);

  // 对输入的图像数据进行文字检测和识别
  char *detectPpocr(Ppocr *Obj, const unsigned char *bitmapData, int width, int height, float nms, float prob, int size, const char *color);

  //
  void closePpocr(Ppocr *Obj);

#ifdef __cplusplus
}
#endif

#endif // PPOCR_H

// // 1. 创建实例
// Ppocr* ocr = newPpocr();

// // 2. 加载模型
// const char* error = loadModelPpocr(ocr, "labels.txt", "db_param", "db_bin",
//                                   "rec_param", "rec_bin", 4);
// if (strlen(error) > 0) {
//     printf("Load model failed: %s\n", error);
//     return;
// }

// // 3. 执行识别（假设有图像数据）
// char* result = detectPpocr(ocr, imageData, 640, 480, 0.3f, 0.5f, 100, "#000000");

// // 4. 处理结果
// printf("OCR Result: %s\n", result);

// // 5. 释放资源
// free(result);  // 释放 detectPpocr 返回的字符串
// closePpocr(ocr);  // 释放 OCR 引擎