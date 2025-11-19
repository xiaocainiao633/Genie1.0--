package agent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// APIDoc API文档结构
type APIDoc struct {
	ID          int    `json:"id"`
	Module      string `json:"module"`
	Function    string `json:"function"`
	Description string `json:"description"`
	Signature   string `json:"signature"`
	Parameters  string `json:"parameters"`
	Return      string `json:"return"`
	Example     string `json:"example"`
	Keywords    string `json:"keywords"` // 用于检索的关键词
}

// KnowledgeBase 知识库
type KnowledgeBase struct {
	db       *sql.DB
	embedder Embedder
}

// Embedder 向量化接口
type Embedder interface {
	Embed(text string) ([]float32, error)
}

// NewKnowledgeBase 创建知识库
func NewKnowledgeBase(dbPath string) (*KnowledgeBase, error) {
	db, err := sql.Open("sqlite3", dbPath + "?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	kb := &KnowledgeBase{db: db}
	if err := kb.initDB(); err != nil {
		return nil, err
	}

	return kb, nil
}

// initDB 初始化数据库
func (kb *KnowledgeBase) initDB() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS api_docs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		module TEXT NOT NULL,
		function TEXT NOT NULL,
		description TEXT,
		signature TEXT,
		parameters TEXT,
		return_type TEXT,
		example TEXT,
		keywords TEXT,
		embedding TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_module ON api_docs(module);
	CREATE INDEX IF NOT EXISTS idx_function ON api_docs(function);
	CREATE INDEX IF NOT EXISTS idx_keywords ON api_docs(keywords);
	`

	_, err := kb.db.Exec(createTableSQL)
	return err
}

// AddAPI 添加API文档
func (kb *KnowledgeBase) AddAPI(doc APIDoc) error {
	insertSQL := `
	INSERT INTO api_docs (module, function, description, signature, parameters, return_type, example, keywords)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := kb.db.Exec(insertSQL,
		doc.Module, doc.Function, doc.Description, doc.Signature,
		doc.Parameters, doc.Return, doc.Example, doc.Keywords, "")
	return err
}

// Search 搜索相关API
func (kb *KnowledgeBase) Search(query string, limit int) ([]APIDoc, error) {
	if limit <= 0 {
		limit = 10
	}

	// 简单的关键词匹配搜索
	queryLower := strings.ToLower(query)
	searchSQL := `
	SELECT id, module, function, description, signature, parameters, return_type, example, keywords
	FROM api_docs
	WHERE 
		LOWER(function) LIKE ? OR
		LOWER(description) LIKE ? OR
		LOWER(keywords) LIKE ? OR
		LOWER(module) LIKE ?
	ORDER BY 
		CASE 
			WHEN LOWER(function) LIKE ? THEN 1
			WHEN LOWER(description) LIKE ? THEN 2
			WHEN LOWER(keywords) LIKE ? THEN 3
			ELSE 4
		END
	LIMIT ?
	`

	pattern := "%" + queryLower + "%"
	rows, err := kb.db.Query(searchSQL, pattern, pattern, pattern, pattern, pattern, pattern, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []APIDoc
	for rows.Next() {
		var doc APIDoc
		err := rows.Scan(&doc.ID, &doc.Module, &doc.Function, &doc.Description,
			&doc.Signature, &doc.Parameters, &doc.Return, &doc.Example, &doc.Keywords)
		if err != nil {
			continue
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// GetByModule 根据模块获取API
func (kb *KnowledgeBase) GetByModule(module string) ([]APIDoc, error) {
	rows, err := kb.db.Query(`
		SELECT id, module, function, description, signature, parameters, return_type, example, keywords
		FROM api_docs
		WHERE module = ?
	`, module)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []APIDoc
	for rows.Next() {
		var doc APIDoc
		err := rows.Scan(&doc.ID, &doc.Module, &doc.Function, &doc.Description,
			&doc.Signature, &doc.Parameters, &doc.Return, &doc.Example, &doc.Keywords)
		if err != nil {
			continue
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// Close 关闭数据库连接
func (kb *KnowledgeBase) Close() error {
	return kb.db.Close()
}

// SetEmbedder 设置Embedding模型
func (kb *KnowledgeBase) SetEmbedder(embedder Embedder) {
	kb.embedder = embedder
}

// EnsureEmbeddings 为知识库生成Embedding
func (kb *KnowledgeBase) EnsureEmbeddings() error {
	if kb.embedder == nil {
		return fmt.Errorf("embedder 未配置，无法生成向量")
	}

	rows, err := kb.db.Query(`SELECT id, description, signature, example FROM api_docs WHERE embedding IS NULL OR embedding = ''`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var desc, sig, example string
		if err := rows.Scan(&id, &desc, &sig, &example); err != nil {
			continue
		}

		text := strings.Join([]string{desc, sig, example}, "\n")
		vector, err := kb.embedder.Embed(text)
		if err != nil {
			return err
		}

		if err := kb.saveEmbedding(id, vector); err != nil {
			return err
		}
	}

	return nil
}

func (kb *KnowledgeBase) saveEmbedding(id int, embedding []float32) error {
	data, err := json.Marshal(embedding)
	if err != nil {
		return err
	}

	_, err = kb.db.Exec(`UPDATE api_docs SET embedding = ? WHERE id = ?`, string(data), id)
	return err
}

func (kb *KnowledgeBase) fetchAllEmbeddings() ([]APIDoc, [][]float32, error) {
	rows, err := kb.db.Query(`SELECT id, module, function, description, signature, parameters, return_type, example, keywords, embedding FROM api_docs WHERE embedding IS NOT NULL AND embedding != ''`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var docs []APIDoc
	var vectors [][]float32

	for rows.Next() {
		var doc APIDoc
		var embeddingStr string
		if err := rows.Scan(&doc.ID, &doc.Module, &doc.Function, &doc.Description,
			&doc.Signature, &doc.Parameters, &doc.Return, &doc.Example, &doc.Keywords, &embeddingStr); err != nil {
			continue
		}

		var vector []float32
		if err := json.Unmarshal([]byte(embeddingStr), &vector); err != nil {
			continue
		}

		docs = append(docs, doc)
		vectors = append(vectors, vector)
	}

	return docs, vectors, nil
}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0
	}

	var dot float64
	var normA float64
	var normB float64

	for i := 0; i < len(a); i++ {
		dot += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	denominator := math.Sqrt(normA) * math.Sqrt(normB)
	if denominator == 0 {
		return 0
	}

	return dot / denominator
}

// BuildDefaultKnowledgeBase 构建默认知识库
func BuildDefaultKnowledgeBase(kb *KnowledgeBase) error {
	fmt.Println("🔧 开始构建默认知识库...")
	apis := []APIDoc{
		// Motion API
		{
			Module:      "motion",
			Function:    "Click",
			Description: "在指定坐标执行点击操作",
			Signature:   "func Click(x, y, fingerID int)",
			Parameters:  "x: X坐标, y: Y坐标, fingerID: 手指ID(1-10)",
			Return:      "无",
			Example:     "motion.Click(100, 200, 1)",
			Keywords:    "点击 触摸 click tap touch 坐标",
		},
		{
			Module:      "motion",
			Function:    "LongClick",
			Description: "在指定坐标执行长按操作",
			Signature:   "func LongClick(x, y, duration int)",
			Parameters:  "x: X坐标, y: Y坐标, duration: 长按时长(毫秒)",
			Return:      "无",
			Example:     "motion.LongClick(100, 200, 500)",
			Keywords:    "长按 long press 按住",
		},
		{
			Module:      "motion",
			Function:    "Swipe",
			Description: "执行滑动操作",
			Signature:   "func Swipe(x1, y1, x2, y2, duration int)",
			Parameters:  "x1,y1: 起始坐标, x2,y2: 结束坐标, duration: 滑动时长",
			Return:      "无",
			Example:     "motion.Swipe(100, 200, 300, 400, 500)",
			Keywords:    "滑动 swipe 拖拽 drag 滑动",
		},
		{
			Module:      "motion",
			Function:    "Back",
			Description: "点击返回键",
			Signature:   "func Back()",
			Parameters:  "无",
			Return:      "无",
			Example:     "motion.Back()",
			Keywords:    "返回 back 后退",
		},
		{
			Module:      "motion",
			Function:    "Home",
			Description: "点击Home键",
			Signature:   "func Home()",
			Parameters:  "无",
			Return:      "无",
			Example:     "motion.Home()",
			Keywords:    "主页 home 首页",
		},

		// UIACC API
		{
			Module:      "uiacc",
			Function:    "New",
			Description: "创建UI控件选择器",
			Signature:   "func New() *Uiacc",
			Parameters:  "无",
			Return:      "*Uiacc: UI选择器对象",
			Example:     "uiacc.New()",
			Keywords:    "选择器 selector ui 控件",
		},
		{
			Module:      "uiacc",
			Function:    "Text",
			Description: "按文本查找控件",
			Signature:   "func (a *Uiacc) Text(value string) *Uiacc",
			Parameters:  "value: 文本内容",
			Return:      "*Uiacc: 链式调用返回选择器",
			Example:     "uiacc.New().Text(\"登录\")",
			Keywords:    "文本 text 文字 查找",
		},
		{
			Module:      "uiacc",
			Function:    "FindOnce",
			Description: "查找单个控件",
			Signature:   "func (a *Uiacc) FindOnce() *UiObject",
			Parameters:  "无",
			Return:      "*UiObject: 找到的控件对象，未找到返回nil",
			Example:     "uiacc.New().Text(\"确定\").FindOnce()",
			Keywords:    "查找 find 搜索 定位",
		},
		{
			Module:      "uiacc",
			Function:    "WaitFor",
			Description: "等待控件出现",
			Signature:   "func (a *Uiacc) WaitFor(timeout int) *UiObject",
			Parameters:  "timeout: 超时时间(毫秒)，0表示无限等待",
			Return:      "*UiObject: 找到的控件对象",
			Example:     "uiacc.New().Id(\"button1\").WaitFor(5000)",
			Keywords:    "等待 wait 超时 timeout",
		},
		{
			Module:      "uiacc",
			Function:    "Click",
			Description: "点击UI控件",
			Signature:   "func (u *UiObject) Click() bool",
			Parameters:  "无",
			Return:      "bool: 是否点击成功",
			Example:     "uiacc.New().Text(\"确定\").FindOnce().Click()",
			Keywords:    "点击 click 控件点击",
		},
		{
			Module:      "uiacc",
			Function:    "SetText",
			Description: "设置输入框文本",
			Signature:   "func (u *UiObject) SetText(str string) bool",
			Parameters:  "str: 要输入的文本",
			Return:      "bool: 是否设置成功",
			Example:     "uiacc.New().Editable(true).FindOnce().SetText(\"Hello\")",
			Keywords:    "输入 input 文本 text 设置",
		},

		// OpenCV API
		{
			Module:      "opencv",
			Function:    "FindImage",
			Description: "在屏幕中查找匹配的图片模板",
			Signature:   "func FindImage(x1, y1, x2, y2 int, template *[]byte, isGray bool, scalingFactor, sim float32) (int, int)",
			Parameters:  "x1,y1: 搜索区域左上角, x2,y2: 搜索区域右下角, template: 模板图片字节, isGray: 是否灰度, scalingFactor: 缩放因子, sim: 相似度",
			Return:      "(int, int): 找到的坐标，未找到返回(-1, -1)",
			Example:     "x, y := opencv.FindImage(0, 0, 0, 0, &templateBytes, false, 1.0, 0.8)",
			Keywords:    "图像 图片 模板 匹配 find image template",
		},

		// PPOCR API
		{
			Module:      "ppocr",
			Function:    "Ocr",
			Description: "在屏幕指定区域进行OCR文字识别",
			Signature:   "func Ocr(x1, y1, x2, y2 int, colorStr string) []Result",
			Parameters:  "x1,y1: 区域左上角, x2,y2: 区域右下角, colorStr: 颜色过滤",
			Return:      "[]Result: 识别结果数组",
			Example:     "results := ppocr.Ocr(0, 0, 1080, 1920, \"\")",
			Keywords:    "OCR 文字识别 识别文字 文本识别",
		},
		{
			Module:      "ppocr",
			Function:    "OcrFromImage",
			Description: "从图像对象进行OCR识别",
			Signature:   "func OcrFromImage(img *image.NRGBA, colorStr string) []Result",
			Parameters:  "img: 图像对象, colorStr: 颜色过滤",
			Return:      "[]Result: 识别结果数组",
			Example:     "results := ppocr.OcrFromImage(img, \"\")",
			Keywords:    "OCR 图像识别",
		},

		// Images API
		{
			Module:      "images",
			Function:    "CaptureScreen",
			Description: "截取屏幕指定区域",
			Signature:   "func CaptureScreen(x1, y1, x2, y2 int) *image.NRGBA",
			Parameters:  "x1,y1: 区域左上角, x2,y2: 区域右下角，0表示全屏",
			Return:      "*image.NRGBA: 图像对象",
			Example:     "img := images.CaptureScreen(0, 0, 0, 0)",
			Keywords:    "截图 屏幕 capture screen",
		},

		// App API
		{
			Module:      "app",
			Function:    "Launch",
			Description: "启动应用",
			Signature:   "func Launch(packageName string, displayId int) bool",
			Parameters:  "packageName: 应用包名, displayId: 显示ID",
			Return:      "bool: 是否启动成功",
			Example:     "app.Launch(\"com.example.app\", 0)",
			Keywords:    "启动 launch 打开 open app",
		},
		{
			Module:      "app",
			Function:    "CurrentPackage",
			Description: "获取当前应用包名",
			Signature:   "func CurrentPackage() string",
			Parameters:  "无",
			Return:      "string: 包名",
			Example:     "pkg := app.CurrentPackage()",
			Keywords:    "包名 package 当前应用",
		},
		{
			Module:      "app",
			Function:    "ForceStop",
			Description: "强制停止应用",
			Signature:   "func ForceStop(packageName string)",
			Parameters:  "packageName: 应用包名",
			Return:      "无",
			Example:     "app.ForceStop(\"com.example.app\")",
			Keywords:    "停止 stop 关闭",
		},

		// IME API
		{
			Module:      "ime",
			Function:    "InputText",
			Description: "输入文本",
			Signature:   "func InputText(text string)",
			Parameters:  "text: 要输入的文本",
			Return:      "无",
			Example:     "ime.InputText(\"Hello World\")",
			Keywords:    "输入 input text 文本",
		},
		{
			Module:      "ime",
			Function:    "SetClipText",
			Description: "设置剪切板文本",
			Signature:   "func SetClipText(text string) bool",
			Parameters:  "text: 文本内容",
			Return:      "bool: 是否成功",
			Example:     "ime.SetClipText(\"Hello\")",
			Keywords:    "剪切板 clipboard",
		},

		// Utils API
		{
			Module:      "utils",
			Function:    "Sleep",
			Description: "等待指定时间",
			Signature:   "func Sleep(i int)",
			Parameters:  "i: 等待时间(毫秒)",
			Return:      "无",
			Example:     "utils.Sleep(1000)",
			Keywords:    "等待 sleep 延时 delay",
		},
	}

	for _, api := range apis {
		if err := kb.AddAPI(api); err != nil {
			return fmt.Errorf("添加API失败 %s.%s: %v", api.Module, api.Function, err)
		}
	}

	return nil
}

// GetContext 获取上下文信息（用于RAG）
func (kb *KnowledgeBase) GetContext(query string) (string, error) {
	var docs []APIDoc
	var err error

	if kb.embedder != nil {
		if err := kb.EnsureEmbeddings(); err != nil {
			return "", err
		}
		docs, err = kb.SearchWithEmbeddings(query, 5)
	} else {
		docs, err = kb.Search(query, 5)
	}

	if err != nil {
		return "", err
	}

	var context strings.Builder
	context.WriteString("相关API文档:\n\n")
	for i, doc := range docs {
		context.WriteString(fmt.Sprintf("%d. %s.%s\n", i+1, doc.Module, doc.Function))
		context.WriteString(fmt.Sprintf("   描述: %s\n", doc.Description))
		context.WriteString(fmt.Sprintf("   签名: %s\n", doc.Signature))
		context.WriteString(fmt.Sprintf("   参数: %s\n", doc.Parameters))
		context.WriteString(fmt.Sprintf("   返回: %s\n", doc.Return))
		context.WriteString(fmt.Sprintf("   示例: %s\n\n", doc.Example))
	}

	return context.String(), nil
}

// SearchWithEmbeddings 使用向量检索相关API
func (kb *KnowledgeBase) SearchWithEmbeddings(query string, limit int) ([]APIDoc, error) {
	if kb.embedder == nil {
		return kb.Search(query, limit)
	}

	vector, err := kb.embedder.Embed(query)
	if err != nil {
		return nil, err
	}

	docs, embeddings, err := kb.fetchAllEmbeddings()
	if err != nil {
		return nil, err
	}

	type scoredDoc struct {
		doc   APIDoc
		score float64
	}

	var scored []scoredDoc
	for i, doc := range docs {
		score := cosineSimilarity(vector, embeddings[i])
		scored = append(scored, scoredDoc{doc: doc, score: score})
	}

	// 排序
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if limit > len(scored) || limit <= 0 {
		limit = len(scored)
	}

	result := make([]APIDoc, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, scored[i].doc)
	}

	return result, nil
}

