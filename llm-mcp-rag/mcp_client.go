package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPClient struct {
	Ctx    context.Context // 上下文
	Client *client.Client  // 底层MCP客户端实例
	Tools  []mcp.Tool      // 缓存从工具进程获取的工具列表
	Cmd    string          // 命令
	Args   []string        // 参数
	Env    []string        // 环境变量
}

// 构造函数
func NewMCPClient(ctx context.Context, cmd string, env, args []string) *MCPClient {
	stdioTransport := transport.NewStdio(cmd, env, args...)
	cli := client.NewClient(stdioTransport)
	m := &MCPClient{
		Ctx:    ctx,
		Client: cli,
		Cmd:    cmd,
		Args:   args,
		Env:    env,
	}
	return m
}

// 启动并初始化MCP连接
func (m *MCPClient) Start() error {
	err := m.Client.Start(m.Ctx)
	if err != nil {
		return err
	}
	mcpInitReq := mcp.InitializeRequest{}
	mcpInitReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	mcpInitReq.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "0.0.1",
	}
	if _, err = m.Client.Initialize(m.Ctx, mcpInitReq); err != nil {
		fmt.Println("mcp init error:", err)
		return err
	}
	return err
}

// 发现并缓存工具列表
func (m *MCPClient) SetTools() error {
	toolsReq := mcp.ListToolsRequest{}
	tools, err := m.Client.ListTools(m.Ctx, toolsReq)
	if err != nil {
		return err
	}
	mt := make([]mcp.Tool, 0)
	for _, tool := range tools.Tools {
		mt = append(mt, mcp.Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}
	m.Tools = mt
	return nil
}

// 调用具体工具
func (m *MCPClient) CallTool(name string, args any) (string, error) {
	var arguments map[string]any
	switch v := args.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &arguments); err != nil {
			return "", err
		}
	case map[string]any:
		arguments = v
	default:
	}
	res, err := m.Client.CallTool(m.Ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      name,
			Arguments: arguments,
		},
	})
	if err != nil {
		return "", err
	}
	return mcp.GetTextFromContent(res.Content), nil
}

// 资源清理
func (m *MCPClient) Close() {
	_ = m.Client.Close()
}

// 获取工具列表
func (m *MCPClient) GetTool() []mcp.Tool {
	return m.TOols
}