package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"

	// ajuste o caminho conforme necessário

	countToken "github.com/gpt_chat/chat_service/internal/utils"
)

type Message struct {
	ID        string
	Role      string
	Content   string
	Tokens    int
	Model     *Model
	CreatedAt time.Time
}

func NewMessage(role, content string, model *Model) (*Message, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    role,
			Content: content,
		},
	}

	totalTokens := countToken.NumTokensFromMessages(messages, model.Name) // Use a função para contar tokens do conteúdo

	msg := &Message{
		ID:        uuid.New().String(),
		Role:      role,
		Content:   content,
		Tokens:    totalTokens,
		Model:     model,
		CreatedAt: time.Now(),
	}
	if err := msg.Validate(); err != nil { // return msg in null if have message erro or return msg
		return nil, err
	}
	return msg, nil
}

func (m *Message) Validate() error {
	if m.Role != "user" && m.Role != "system" && m.Role != "assistant" {
		return errors.New("invalid role")
	}
	if m.Content == "" {
		return errors.New("content is empty")
	}
	if m.CreatedAt.IsZero() {
		return errors.New("invalid created at")
	}
	if m.Tokens == 0 {
		return errors.New("token is empty")
	}
	return nil
}

func (m *Message) getQtdTokens() int {
	return m.Tokens
}
