package agent

import (
	"fmt"
	"strings"
	"time"
)

// Message 对话消息
type Message struct {
	Role      string
	Content   string
	Timestamp time.Time
}

// ConversationMemory 多轮对话记忆
type ConversationMemory struct {
	history []Message
	limit   int
}

// NewConversationMemory 创建记忆
func NewConversationMemory(limit int) *ConversationMemory {
	if limit <= 0 {
		limit = 5
	}
	return &ConversationMemory{
		history: make([]Message, 0, limit),
		limit:   limit,
	}
}

// AddMessage 添加记录
func (m *ConversationMemory) AddMessage(role, content string) {
	if strings.TrimSpace(content) == "" {
		return
	}
	m.history = append(m.history, Message{
		Role:      role,
		Content:   strings.TrimSpace(content),
		Timestamp: time.Now(),
	})
	if len(m.history) > m.limit {
		m.history = m.history[len(m.history)-m.limit:]
	}
}

// ContextString 返回上下文
func (m *ConversationMemory) ContextString() string {
	if len(m.history) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString("对话上下文:\n")
	for _, msg := range m.history {
		builder.WriteString(fmt.Sprintf("[%s] %s\n", msg.Role, msg.Content))
	}
	return builder.String()
}

// LastUserQuery 返回最近的用户消息
func (m *ConversationMemory) LastUserQuery() string {
	for i := len(m.history) - 1; i >= 0; i-- {
		if m.history[i].Role == "user" {
			return m.history[i].Content
		}
	}
	return ""
}

