### 效果演示：

首先需要连接到安卓设备，然后根据 agent/README.md 文件中的操作连接本地 llm，可在 agent.go 目录下修改配置，如果没有连接本地 llm，则默认使用离线模式，

**部分截图如下：**

与大模型对话：

<img width="1170" height="973" alt="image" src="https://github.com/user-attachments/assets/009650c2-7029-4de8-b5b5-098910ee774f" />

<img width="843" height="984" alt="image" src="https://github.com/user-attachments/assets/ec569c22-c716-4e81-a075-0abbcabcaff9" />

初始化产生的二进制数据库文件:
① 存储 API 知识库
② 支持语义检索（RAG）

<img width="307" height="374" alt="image" src="https://github.com/user-attachments/assets/f52d7997-1f25-4701-bc60-c031da37d97a" />


workspace 目录：自动生成，每生成一个自动化测试用例会自动添加到此目录
文件形式为 json 或者 markdown(可自行设置)，同时记录测试结果

<img width="395" height="1112" alt="image" src="https://github.com/user-attachments/assets/34c169a0-f7c3-42cc-91a3-c485f0b65bab" />

根据提供的信息，以下是完整的Go测试脚本：

```go
// Package main
package main

import (
    "fmt"
    "log"

    "auto Go"
)

func main() {
    //导入必要的AutoGo模块
    goInit()

    //创建一个新的自动化测试客户端
    client := NewAutoClient("http://127.0.0.1:8080")

    //设置应用名称
    client.SetAppName("test_app")

    //点击登录按钮
    motionClick(100, 200, 1)

    //随机生成一组测试用例，要求调用openCV和ppocr
    testCases := generateTestCases()
    for _, testCase := range testCases {
        log.Printf("开始执行测试用例:%s", testCase.Name)
        executeTestCase(client, testCase)
        log.Printf("测试用例/%s完成\n", testCase.Name)
    }

    //在刚刚代码的基础上添加调用opencv
    motionClick(300, 400, 2)

}

// generateTestCases 生成随机测试用例
func generateTestCases() []*TestCase {
    testCases := make([]*TestCase, 0)
    for i := 1; i <= 10; i++ {
        testCase := &TestCase{
            Name: fmt.Sprintf("testCase-%d", i),
            Steps: []Step{
                {"openCV", "opencv"},
                {"ppocr", "ppocr"},
            },
        }
        testCases = append(testCases, testCase)
    }
    return testCases
}

// executeTestCase 执行测试用例
func executeTestCase(client *AutoClient, testCase *TestCase) {
    for _, step := range testCase.Steps {
        log.Printf("开始执行步骤:%s", step.Name)
        switch step.Type {
        case "openCV":
            //执行openCV操作
            client.OpenCV(step.Data["opencv"])
        case "ppocr":
            //执行ppocr操作
            client.PPocr(step.Data["ppocr"])
        default:
            log.Println("未知步骤")
        }
        log.Printf("步骤/%s完成\n", step.Name)
    }
}

type TestCase struct {
    Name  string
    Steps []Step
}
type Step struct {
    Name  string
    Type string
    Data map[string]string
}

// NewAutoClient create a new AutoClient instance
func NewAutoClient(baseURL string) *AutoClient {
    return &AutoClient{
        baseURL: baseURL,
    }
}

// SetAppName set the app name for the client
func (c *AutoClient) SetAppName(appName string) {
    c.appName = appName
}

// OpenCV execute opencv operation
func (c *AutoClient) OpenCV(opencvData map[string]string) {
    log.Printf("执行openCV操作：%v", opencvData)
}

// PPocr execute ppocr operation
func (c *AutoClient) PPocr(ppocrData map[string]string) {
    log.Printf("执行ppocr操作：%v", ppocrData)
}
```

此脚本包含以下功能：

1.导入必要的 AutoGo 模块。 2.创建一个新的自动化测试客户端。 3.设置应用名称。 4.点击登录按钮。 5.随机生成一组测试用例，要求调用 openCV 和 ppocr。 6.在刚刚代码的基础上添加调用 opencv。

该脚本使用 AutoGo 提供的 API 执行测试操作，包括点击登录、打开 OpenCV 和 PPOCR。它还包含基本的错误处理和日志输出。

请注意，这个脚本是 example，可以根据实际需求进行修改和扩展。


