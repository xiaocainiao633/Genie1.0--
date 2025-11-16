#ifndef IMGUI_H
#define IMGUI_H

#ifdef __cplusplus
extern "C"
{
#endif

  int Init(int hide);
  int CheckInit();
  void Toast(const char *message);
  void Toast_setTextSize(int size);
  void Console_println(const char *message);
  void Console_setSize(int w, int h);
  void Console_setPosition(int x, int y);
  void Console_setTextColor(const char *color);
  void Console_setWindowColor(const char *color);
  void Console_setTextSize(int size);
  void Console_clear();
  void Console_hide();
  void Console_show();
  void Rect_add(int x1, int y1, int x2, int y2, const char *color);
  void Rect_clear();
  void StrLine_add(int x1, int y1, int x2, int y2, const char *color);
  void StrLine_clear();
  void Image_add(int x1, int y1, int x2, int y2, unsigned char *data, int dataLeng);
  void Image_clear();
  void Hud_init(int x1, int y1, int x2, int y2, const char *color, int textSize);
  void Hud_setText(const char *str);
  void Hud_clear();
  void Alert(const char *title, const char *message);
  void Close();

#ifdef __cplusplus
}
#endif

#endif // IMGUI_H
