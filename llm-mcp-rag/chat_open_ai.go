package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

// 存储必要信息
type ChatOpenAI struct {
	Ctx						context.Context  // 上下文
	Model 				string           // 要使用的LLM模型
	SystemPrompt 	string           // 系统提示词
	Tools 				[]mcp.Tool       // MCP工具集
	RagContext  	string           // RAG上下文  
	Message 			[]openai.ChatCompletionMessageParamUnion  // 聊天历史信息
	LLM 					openai.Client    // 官方客户端实例
}

// 专门用于修改ChatOpenAI的属性
type LLMOption func(*ChatOpenAI)

// 设置系统提示词
func WithSystemPrompt(prompt string) LLMOption {
	return func (ai *ChatOpenAI) {
		ai.SystemPrompt = prompt  // 把用户传入的prmpt赋值给ChatOpenAI
	}
}

func WithRagContext(ragPrompt string) LLMOption {
	return func(ai *ChatOpenAI) {
		ai.RagContext = ragPrompt // 添加 RAG 上下文
	}
}

func WithMessage(message []openai.ChatCompletionMessageParamUnion) LLMOption {
	return func(ai *ChatOpenAI) {
		ai.Message = message // 添加聊天信息
	}
}

func WithLLMTools(tools []mcp.Tool) LLMOption {
	return func(ai *ChatOpenAI) {
		ai.Tools = tools // 添加工具
	}
}

// 从环境变量读取apikey和base_url等，构建openai客户端
func NewChatOpenAI(ctx context.Context, model string, opts ...LLMOption) *ChatOpenAI {
	if model == "" {
		panic("model is required")
	}
	var (
		apiKey = os.Getenv(ChatGPTOpenAPIKEY)
		baseURL = os.Getenv(ChatGPTBaseURL)
	)
	if apiKey == "" {
		panic("missing OPENAI_API_KEY")
	}
	options := []option.RequestOption{
		option.WithAPIKEY(apiKey)
	}
	if baseURL != "" {
		options = append(options, option.WithBaseUrl(baseurl))
	}
	cli := openai.NewClient(options...)
	llm := &ChatOpenAI{
		Ctx:     ctx,
		Model:   model,
		LLM:     cli,
		Message: make([]openai.ChatCompletionMessageParamUnion, 0),
	}
	for _, opt := range opts {
		opt(llm)
	}
	if llm.SystemPrompt != "" {
		llm.Message = append(llm.Message, openai.SystemMessage(llm.SystemPrompt))
	}
	if llm.RagContext != "" {
		llm.Message = append(llm.Message, openai.UserMessage(llm.RagContext))
	}
	fmt.Println("init LLM successfully")
	return llm
}

// 核心对话逻辑
// 输入用户当前输入文本，输出自然语言回复和模型请求调用的工具列表
func (c *ChatOpenAI) Chat(prompt string) (result string, toolCall []openai.ToolCallUnion) {
	fmt.Println("init chat...")
	// 如果输入非空，将其加入到对话历史中
	if prompt != "" {
		// 追加用户消息
		c.Message = append(c.Message, openai.UserMessage(prompt))
	}
	// 将内部存储的MCP格式工具列表转换为OpenAI SDK要求的工具参数格式
	toolsParam := MCPTool2OpenAITool(c.Tools)
	if len(toolsParam) == 0 {
		toolsParam = nil
	}
	// 调用聊天接口
	stream := c.LLM.Chat.Completions.NewStreaming(c.Ctx, openai.ChatCompletionNewParams{
		Messages: c.Message,
		Seed:     openai.Int(0),
		Model:    c.Model,
		Tools:    toolsParam,
	})
	acc := openai.ChatCompletionAccumulator{}
	var toolCalls []openai.ToolCallUnion
	result = ""
	finished := false
	fmt.Println("start chatting...")
	// 迭代流式响应
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)
		if content, ok := acc.JustFinishedContent(); ok {
			finished = true
			result = content
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Println("tool call finished:", tool.Index, tool.Name, tool.Arguments)
			toolCalls = append(toolCalls, openai.ToolCallUnion{
				ID: tool.ID,
				Function: openai.FunctionToolCallFunction{
					Name:      tool.Name,
					Arguments: tool.Arguments,
				},
			})
		}
		if refusal, ok := acc.JustFinishedRefusal(); ok {
			fmt.Println("refusal:", refusal)
		}

		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta.Content
			if !finished {
				result += delta
			}
		}
	}

	if len(acc.Choices) > 0 {
		c.Message = append(c.Message, acc.Choices[0].Message.ToParam())
	}

	if stream.Err() != nil {
		panic(stream.Err())
	}

	return result, toolCalls
}

// 将MCP工具描述转换为能识别的工具定义
func MCPTool2OpenAITool(mcpTools []mcp.Tool) []openai.ChatCompletionToolUnionParam {
	// 存储转换后的OpenAI工具定义
	openAITools := make([]openai.ChatCompletionToolUnionParam, 0)
	for _, tool := range mcpTools {
		// 构建响应风格的参数描述
		params := openai.FunctionParameters{
			"type":       tool.InputSchema.Type,
			"properties": tool.InputSchema.Properties,
			"required":   tool.InputSchema.Required,
		}
		// 关键兜底：若type为空，默认用object，避免OpenAI拒绝工具定义
		if t, ok := params["type"].(string); !ok || t == "" {
			params["type"] = "object"
		}
		openAITools = append(openAITools, openai.ChatCompletionToolUnionParam{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: shared.FunctionDefinitionParam{
					Name:        tool.Name,
					Description: openai.String(tool.Description),
					Parameters:  params,
				},
			},
		})
	}
	return openAITools
}