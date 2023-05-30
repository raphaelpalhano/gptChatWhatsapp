package gateway

import (
	"context"

	"github.com/gpt_chat/chat_service/internal/domain/entity"
)

type ChatGatway interface {
	CreateChat(ctx context.Context, chat *entity.Chat) error
	FindChatById(ctx context.Context, chatID string) (*entity.Chat, error)
	SaveChat(ctx context.Context, chat *entity.Chat) error
}
