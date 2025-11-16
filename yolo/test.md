**在目标区域内根据坐标识别特定物体,返回标签以及坐标**

## UI检测

```
func TestUIComponents() {
    yolo := yolo.New("v5", 2, "ui_elements.param", "ui_elements.bin", "ui_labels.txt")
    defer yolo.Close()
    
    // 检测特定区域的UI元素
    results := yolo.Detect(100, 200, 500, 600)
    
    foundButton := false
    foundInput := false
    
    for _, result := range results {
        switch result.Label {
        case "button":
            foundButton = true
            gui.DrawRect(result.X, result.Y, result.X+result.Width, result.Y+result.Height, "#00FF00")
            gui.Toast("找到按钮")
            
        case "input_field":
            foundInput = true  
            gui.DrawRect(result.X, result.Y, result.X+result.Width, result.Y+result.Height, "#0000FF")
            gui.Toast("找到输入框")
        }
    }
    
    // 测试断言
    if !foundButton {
        gui.Alert("测试失败", "未找到按钮元素")
    }
    if !foundInput {
        gui.Alert("测试失败", "未找到输入框元素")
    }
}
```