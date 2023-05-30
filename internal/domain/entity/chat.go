package entity

import (
	"errors"

	"github.com/google/uuid"
)

// Config chat

type ChatConfig struct {
	Model            *Model
	Temperature      float32  // precision of answers
	TopP             float32  // Definition of choice the words
	N                int      // number of messages generate
	Stop             []string // stop chat
	MaxTokens        int      // number max of tokens
	PresencePenalty  float32  // anal
	FrequencyPenalty float32
}

type Chat struct {
	ID                   string
	UserID               string
	InitialSystemMessage *Message
	Messages             []*Message
	ErasedMessages       []*Message // messages deleted for create space for more tokens
	Status               string
	TokenUsage           int
	Config               *ChatConfig
}

/*
Add new message chat:

1. If message is more bigger than maxTokens -->
Remove before messages

algorith: run of all messages if exceed return error
*/

func NewChat(userID string, initialSystemMessage *Message, chatConfig *ChatConfig) (*Chat, error) {

	chat := &Chat{
		ID:                   uuid.New().String(),
		UserID:               userID,
		InitialSystemMessage: initialSystemMessage,
		Status:               "active",
		Config:               chatConfig,
		TokenUsage:           0,
	}

	chat.AddMessage(initialSystemMessage)

	if err := chat.Validate(); err != nil {
		return nil, err
	}

	return chat, nil
}

func (c *Chat) Validate() error {
	if c.UserID == "" {
		return errors.New("user id is empty")
	}
	if c.Status != "active" && c.Status != "ended" {
		return errors.New("invalid status")
	}
	if c.Config.Temperature < 0 || c.Config.Temperature > 2 {
		return errors.New("invalid temperature")
	}
	// ...
	return nil
}

func (c *Chat) AddMessage(m *Message) error {
	if c.Status == "ended" {
		return errors.New("Chat is ended. No more messages allowed")
	}
	for {
		if c.Config.Model.GetMaxTokens() >= m.getQtdTokens()+c.TokenUsage {
			c.Messages = append(c.Messages, m)
			c.RefreshTokenUsage()
			break
		}
		c.ErasedMessages = append(c.ErasedMessages, c.Messages[0])
		c.Messages = c.Messages[1:]
		c.RefreshTokenUsage()
	}
	return nil
}

func (c *Chat) GetMessages() []*Message {
	return c.Messages
}

func (c *Chat) CountMessage() int {
	return len(c.Messages)
}

func (c *Chat) End() error {
	if c.Status == "ended" {
		return errors.New("Chat already ended")
	}
	c.Status = "ended"

	return nil

}

func (c *Chat) RefreshTokenUsage() {
	c.TokenUsage = 0
	for m := range c.Messages {
		c.TokenUsage += c.Messages[m].getQtdTokens()

	}
}
