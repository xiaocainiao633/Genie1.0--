# 对于 Genie 的一些分析

原项目仓库地址：[Genie](https://github.com/functional-fuzzing-android-apps/home.git)
新项目地址：[Genie-1.0](https://github.com/xiaocainiao633/Genie1.0--.git)

## 目录

- [1.原项目目录树](#原项目目录树)
- [2.项目文件分析](#项目文件分析)
  - [2.1 app-ActivityDiary-bug-report-example/](#目录一)
  - [2.2 结论一](#结论一)
  - [2.3 Genie/](#目录二)
  - [2.4 结论二](#结论二)
- [3.项目改进](#项目改进)

### 原项目目录树

```
app-ActivityDiary-bug-report-example/
├── events/
├── seed-tests/
│   └── seed-test-7/
├── states/
├── stylesheets/
├── views/
├── index_cluster.html
├── index.html
├── utg_cluster.js
├── utg_cluster.json
├── utg.js
├── utg.json
Genie/
├── apps_for_test/
├── deploy/
├── droidbot/
├── droidbot.egg-info/
├── script_samples/
├── tmp-diary/
├── .gitignore
├── LICENSE
├── setup.cfg
├── setup.py
```

### 项目文件分析

本项目主要使用了 python 和 javascript，先简单阅读一下项目源码：

#### 目录一

包含两个 html 文件：UTG 模型的可视化页面入口，有助于了解模型
events、states、views 中存放的是测试中的状态
.json 文件：可以对照页面上的状态和转移，找到数据与可视化的对应关系，某个节点的详细信息
json 中数据基本可以分为以下两种格式：

```
方式一：
{
      "from": "db12386906c74798d403e64ff361d346",
      "to": "1dacdc6fd3b1c4539eb196feedb2636c",
      "id": "db12386906c74798d403e64ff361d346-->1dacdc6fd3b1c4539eb196feedb2636c",
      "title": "<table class=\"table\">\n<tr><th>6</th><td>TouchEvent(state=db12386906c74798d403e64ff361d346, view=ebc5459783db08890df90764b0a11ada_5)</td></tr>\n</table>",
      "label": "6",
      "events": [
        {
          "event_str": "TouchEvent(state=db12386906c74798d403e64ff361d346, view=ebc5459783db08890df90764b0a11ada_5)",
          "event_id": 6,
          "event_type": "touch",
          "event_log_file_path": "events/event_2020-02-21_220710-7.json",
          "view_images": [
            "views/view_ebc5459783db08890df90764b0a11ada_5_2020-02-21_220712_5.png"
          ]
        }
      ]
    },
```

和

```
方式二：
{
      "id": "037b4d8604fa58e0e5f458b2ce242d9c",
      "shape": "image",
      "image": "states/screen_2020-02-22_000213-91.png",
      "label": "HistoryDetailActivity",
      "package": "de.rampro.activitydiary.debug",
      "activity": "de.rampro.activitydiary.ui.history.HistoryDetailActivity",
      "state_str": "037b4d8604fa58e0e5f458b2ce242d9c",
      "structure_str": "037b4d8604fa58e0e5f458b2ce242d9c",
      "title": "<table class=\"table\">\n<tr><th>package</th><td>de.rampro.activitydiary.debug</td></tr>\n<tr><th>activity</th><td>de.rampro.activitydiary.ui.history.HistoryDetailActivity</td></tr>\n<tr><th>state_str</th><td>40a3cfee4e5a4c548d5b2fdb60cdf2c7</td></tr>\n<tr><th>structure_str</th><td>037b4d8604fa58e0e5f458b2ce242d9c</td></tr>\n</table>",
      "content": "de.rampro.activitydiary.debug\nde.rampro.activitydiary.ui.history.HistoryDetailActivity\n40a3cfee4e5a4c548d5b2fdb60cdf2c7\nandroid:id/content,android:id/animator,android:id/datePicker,android:id/date_picker_header,android:id/button1,android:id/date_picker_header_date,android:id/date_picker_day_picker,android:id/customPanel,android:id/next,android:id/custom,android:id/button2,android:id/date_picker_header_year,android:id/buttonPanel,android:id/parentPanel,android:id/day_picker_view_pager,android:id/prev,android:id/month_view\n2,11,5,29,16,1,26,3,14,17,24,22,9,Sat, Feb 22,2020,23,6,8,12,OK,21,20,18,19,30,10,Cancel,15,25,7,28,13,4,31,27"
    }
```

通过这分析两种方式，我们大致可以明白测试的基本流程如下：
**1.UI 状态与事件采集阶段：**
此阶段通过自动化工具模拟用户操作，实时捕获应用的界面状态和操作事件，为后续构建 UTG 模型提供原始数据。
自动化操作执行：Genie 目录中的 `DroidBot` 按照预设规则模拟用户行为，如点击按钮、输入文本、滑动界面等。
每执行一次操作后，工具会：
截取当前界面截图，保存到 states/ 目录（如 screen_2020-02-22_000213-91.png）；

**提取界面核心信息（包名、Activity 类名、控件 ID 列表、文本内容），生成 state_str（状态哈希）、structure_str（结构哈希），并封装成 “状态节点” 数据（对应 utg.json 中的 id、activity、content 等字段）。**

同时记录每次操作的详细信息，包括：
操作类型（如 touch 触摸）、操作序号（如 event_id: 6）；
被操作控件的 ID（如 ebc5459783db08890df90764b0a11ada_5）及控件截图（保存到 views/ 目录）；
操作的时间戳、坐标等原始日志，写入 events/ 目录下的 JSON 文件（如 event_2020-02-21_220710-7.json）。

**2. UTG 模型构建阶段**
此阶段将采集到的 “状态节点” 和 “事件数据” 关联，生成完整的 UI 状态转移图（UTG），即 utg.json 和 utg_cluster.json 中的核心数据。
转移关系生成：根据操作顺序，将 “前一状态 → 操作事件 → 后一状态” 关联，形成 “转移边” 数据：
也就是方式一中的数据写法：通过 from 和 to 描述状态之间的跳转关系
前一状态 ID 作为 from，后一状态 ID 作为 to；
操作事件封装到 events 数组（包含 event_str、event_log_file_path 等）；
生成转移边的唯一 ID（fromID-->toID）和简短标识（label）。
模型聚合：若执行多组测试用例，会将多组 UTG 数据合并，生成 utg_cluster.json，形成全局状态转移图（包含所有测试用例覆盖的状态和转移路径）。

**3. 测试验证与缺陷检测阶段**
此阶段基于 UTG 模型，通过 “预期 vs 实际” 的对比，验证应用行为是否符合预期，识别异常或缺陷。
路径覆盖校验：检查 UTG 模型是否覆盖应用的核心功能路径（如 “首页 → 历史列表 → 历史详情”），若存在未覆盖的关键路径，判定为 “测试覆盖不充分”。
状态一致性校验：对比相同操作下的 “预期状态” 与 “实际状态”：
若两次相同操作后，state_str 或 content 不一致（如界面控件缺失、文本错误），判定为 “状态异常”；
若操作后未按预期跳转到目标状态（如 to 与模型中预设的目标 ID 不符），判定为 “跳转异常”。
崩溃与无响应检测：若操作后应用崩溃（Activity 异常退出）或无响应，工具会记录崩溃时的状态节点（截图、日志），并在 UTG 模型中标记该节点为 “缺陷关联状态”，关联对应的操作事件日志。

**4. 缺陷定位与报告生成阶段**
此阶段针对检测到的异常，结合 UTG 模型和日志数据，定位问题根源，并生成可视化报告（index.html 和 index_cluster.html）。
此时，我们可以查看以下目录中的部分数据
查看 events/ 下对应的操作日志，确认触发异常的具体操作（如触摸坐标、控件 ID）；
查看 views/ 下的控件截图，确认是否操作了错误控件；
查看 states/ 下的界面截图，观察异常时的界面表现（如空白、弹窗错误）。
点击状态节点可查看 title 中的详细属性（包名、Activity 等）；(方式二)
点击转移边可查看 events 中的操作日志和控件截图；(方式一)

#### stylesheets/文件夹：

该文件夹主要是前端展示页面的设计，也就是刚刚 html 页面的样式，包括几个 css 和 js 文件，对数据进行可视化分析

**1. 数据来源于自动化测试：**
①UTG 数据：来自 utg.json（单组状态转移图）或 utg_cluster.json（多组聚合图），包含所有 UI 状态（节点）和转移关系（边）的结构化描述。
② 测试元数据：来自 checking_result.json（多组 UTG 差异）、logcat.txt（系统日志）、coverage.ec（代码覆盖率）等，是测试结果的 “补充说明”。
(在 seed-tests/目录下)

**2. droidbotUI.js：单组 UTG 的展示：**
它设计单组 UI 交互逻辑，核心作用就是将 utg.json 转化为交互式图形。
①UTG 图渲染：通过 vis.js 将 utg.json 中的 “节点（UI 状态）” 和 “边（转移关系）” 渲染成可拖拽、缩放的网络图形。
② 节点：以界面截图为载体，点击可查看状态详情（包名、Activity 类名等）。
③ 边：以带箭头的连线为载体，点击可查看转移事件详情（操作类型、控件截图等）。
④ 测试概览报告：通过 getOverallResult() 生成表格，展示应用、设备、测试的关键指标（如事件数、Activity 覆盖率），让测试规模和效果 “一目了然”。
⑤ 聚类与搜索：通过 clusterStructures()/clusterActivities() 按界面结构或 Activity 分组节点，通过 searchUTG() 按关键词筛选节点。

**3. aligned_view.js：多组 UTG 的对比**
这是缺陷识别的核心对比工具，通过 Vue 实现种子测试（无缺陷）与变异测试（注入缺陷）的 UTG 对齐。
① 将多组 UTG 的状态节点和转移边在数组中一一对应，即使某组 UTG 新增了操作（如变异测试注入了缺陷操作），也会通过插入 null 保证位置同步。
② 从 checking_result.json 中读取 UTG 差异点（delta），在对比表格中用特殊样式标记，直观展示缺陷导致的 UI 变化。

**4. stylesheets 依赖库：前端展示的主要支撑**
①Bootstrap + jQuery：负责页面的响应式布局（如 UTG 图与详情面板的分栏）、交互组件（如弹窗、表格、进度条）。
②Vue.js：是 aligned_view.js 的运行时依赖，通过组件化和数据绑定实现多组 UTG 对比的动态渲染，保证差异展示的实时性。
③vis.js：droidbotUI.js 通过网络图形算法实现 UTG 节点和边的高效渲染，支持大规模图的流畅交互拖拽、缩放、选中）。
④droidbotUI.css/droidbotUI.js：是与 DroidBot 工具的定制化联动层，负责将自动化操作事件（如触摸、输入）与 UTG 图的边关联，让工具行为与可视化结果深度绑定。

#### seed-tests/文件夹：

此时再来看里面的文件，发现基本都是之前分析过的，可以查看一下其中的 html 文件，并且数据量也较少，帮助我们快速定位缺陷导致的 UI 交互变化

### 结论一：

那么验证测试的逻辑就是：数据采集 → 模型构建 → 对比验证 → 可视化报告
**测试执行时，工具会捕获三类核心数据：**
① 界面状态数据：每个界面的截图（states/ 目录）、界面属性（包名、Activity 类名、控件 ID 等）→ 最终转化为 utg.json 中的状态节点。
② 操作事件数据：每次触摸、输入的操作详情（坐标、时间、控件 ID 等）→ 最终转化为 utg.json 中的 ** 转移边** 的 events 字段。
③ 系统日志数据

**采集的数据会被整理成 UI 状态转移图的 JSON 模型（utg.json）：**
① 每个状态节点对应一个界面，包含截图路径、界面属性、唯一 ID 等信息。
② 每个转移边对应一次操作，包含操作类型、触发的控件截图、关联的起始状态 → 目标状态等信息。

**对比逻辑分为两类：**
① 单组 UTG 自检：将实际生成的 UTG 与应用的需求逻辑对比（如预期有添加日记 → 保存成功的转移，若实际没有则判定为缺陷）。
② 多组 UTG 互检：将基准测试（无缺陷）的 UTG（如 utg_8.json）与变异测试（注入缺陷）的 UTG（如 utg_9.json）对比，通过 checking_result.json 标记差异点（如新增的异常状态、缺失的正常转移）。

最终通过前端脚本（droidbotUI.js、aligned_view.js 等）将 JSON 模型和差异数据转化为交互页面：
① 单组 UTG 图（index.html）：展示应用的完整 UI 交互逻辑，点击节点 / 边可查看详情。
② 多组 UTG 对比图（index_aligned_xx.html）：高亮差异点，直观展示缺陷对 UI 交互的影响。
至此，javascript 部分结束。

### 目录二

**文件内容如表所示**
| 目录 / 文件 | 类型 | 核心作用 |
| ------------------ | ------------- | ----------------------------------------------------------------------------------------- |
| apps_for_test/ | 应用存储 | 存放待测试的安卓应用安装包（APK），是自动化测试的基础。 |
| deploy/ | 部署脚本 | 实现测试环境的自动化部署（如将被测应用、测试框架安装到安卓设备 / 模拟器）。 |
| droidbot/ | 自动化工具 | 集成开源安卓自动化测试工具 DroidBot，负责生成 UI 操作事件（点击、输入等），模拟用户行为。 |
| droidbot.egg-info/ | Python 元数据 | 存储 droidbot 工具的 Python 包依赖、版本信息.这是本地配置好之后生成的文件 |
| script_samples/ | 脚本示例 | 提供自动化测试脚本模板，展示如何编写 “UI 操作流程、断言校验、结果分析” 的逻辑。 |
| tmp-diary/ | 临时目录 | 存储测试过程中的临时文件（如日志、截图、中间数据），测试结束后可清理,同样是本地执行测试之后才有的输出 |
| setup.py/setup.cfg | 安装配置 | Python 项目的安装入口，用于安装 Genie 框架的依赖和模块（如 pip install .）。 |

由表格内容可知:只需简单看一下 droidbot/ 和 deploy/ 文件的内容,其他文件是作为辅助测试
在这一步，如果是 windows 系统，那么需要修改一些配置才能正常部署项目，并且按需修改启动等指令,否则无法进行测试

### 结论二

项目的核心是将安卓应用的 UI 交互逻辑转化为状态节点 + 转移边的结构化模型（UTG），通过模型预期 vs 实际执行的对比，识别 UI 缺陷。
完整流程:

1. 测试准备阶段
   被测应用准备：将待测试的 APK 包放入 Genie/apps_for_test/ 目录。
   环境部署：通过 Genie/deploy/ 脚本，确保安卓设备 / 模拟器已连接、依赖已安装。
   UTG 模型初始化：在 app-ActivityDiary-bug-report-example/ 中维护 utg.json（预期 UI 交互模型），作为测试基础
2. 自动化测试执行阶段（核心：Genie/droidbot/）
   应用启动与操作生成：DroidBot 工具自动安装并启动被测应用，基于动态分析生成触摸、输入等操作事件，模拟用户行为。
   界面状态：截取当前界面截图（存入 app-ActivityDiary-bug-report-example/states/），提取包名、Activity 类名、控件结构等信息。
   操作事件：记录每次操作的类型、控件 ID、时间戳等（存入 app-ActivityDiary-bug-report-example/events/）。
   系统日志：捕获安卓设备的运行日志（logcat.txt）、代码覆盖率（coverage.ec）。
   UTG 模型构建：将采集的状态 + 事件整理成 utg.json（实际执行的 UI 交互模型），包含所有状态节点和转移边。
3. 测试验证与缺陷分析阶段
   多组 UTG 对比：若执行了基准测试（无缺陷）和变异测试（注入缺陷），则通过 seed-tests/seed-test-7/mutant-659/checking_result.json 标记两组 utg.json 的差异点（如新增状态、异常转移）。
   单组 UTG 自检：对比实际 utg.json 与预期 utg.json，识别未覆盖路径、异常状态等缺陷。
4. 可视化报告生成阶段
   单组 UTG 可视化：通过 app-ActivityDiary-bug-report-example/index.html（droidbotUI.js 驱动），将 utg.json 渲染成交互式图形，点击节点 / 边可查看详情。
   多组 UTG 对比可视化：通过 seed-tests/seed-test-7/mutant-659/index_aligned_xx.html（aligned_view.js 驱动），高亮展示两组 UTG 的差异点，直观识别缺陷。

### 项目改进

1.首先是对一些命令进行修改，可以在 windows 系统上直接部署，具体命令和原文档一致即可。

2.对 droidbot 进行优化：什么是 droidbot 呢？下面是官方介绍：
DroidBot 是一个轻量级的 Android 测试输入生成器，它能够向 Android 应用程序发送随机或脚本输入事件，更快地实现更高的测试覆盖率，并在测试后生成 UI 转换图（UTG）。DroidBot 的主要优势包括：
不需要系统修改或应用插桩。
事件基于 GUI 模型（而不是随机）。
可编程（可以自定义某些 UI 的输入）。
能够生成 UI 结构和方法跟踪以进行分析。
对于 api 的书写以及使用方式进行了改进，具体改进见仓库的简介
