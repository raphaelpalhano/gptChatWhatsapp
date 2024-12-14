package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gpt_chat/chat_service/internal/domain/entity"
	"github.com/gpt_chat/chat_service/internal/infra/db"
)

type ChatRepositoryMySQL struct {
	DB      *sql.DB
	Queries *db.Queries
}

func NewChatRepositoryMySQL(dbt *sql.DB) *ChatRepositoryMySQL {
	return &ChatRepositoryMySQL{
		DB:      dbt,
		Queries: db.New(dbt),
	}
}

func (r *ChatRepositoryMySQL) CreateChat(ctx context.Context, chat *entity.Chat) error {
	err := r.Queries.CreateChat(
		ctx,
		db.CreateChatParams{
			ID:               chat.ID,
			UserID:           chat.UserID,                           // Assuming chat.UserID exists
			InitialMessageID: chat.InitialSystemMessage.ID,          // Add this if it exists in the chat struct
			Status:           chat.Status,                           // Add this if it exists in the chat struct
			TokenUsage:       int32(chat.TokenUsage),                // Add this if it exists in the chat struct
			Model:            chat.Config.Model.Name,                // Add this if it exists in the chat struct
			ModelMaxTokens:   int32(chat.Config.MaxTokens),          // Add this if it exists in the chat struct
			Temperature:      float64(chat.Config.Temperature),      // Add this if it exists in the chat struct
			TopP:             float64(chat.Config.TopP),             // Add this if it exists in the chat struct
			N:                int32(chat.Config.N),                  // Add this if it exists in the chat struct
			Stop:             chat.Config.Stop[0],                   // Add this if it exists in the chat struct
			MaxTokens:        int32(chat.Config.MaxTokens),          // Add this if it exists in the chat struct
			PresencePenalty:  float64(chat.Config.PresencePenalty),  // Add this if it exists in the chat struct
			FrequencyPenalty: float64(chat.Config.FrequencyPenalty), // Add this if it exists in the chat struct
			CreatedAt:        time.Now(),                            // Set to current time or appropriate value
			UpdatedAt:        time.Now(),                            // Set to current time or appropriate value
		},
	)

	if err != nil {
		return err
	}

	err = r.Queries.AddMessage(
		ctx,
		db.AddMessageParams{
			ID:        chat.ID,                             // Adicione um ID único para a mensagem
			ChatID:    chat.ID,                             // ID do chat associado
			Role:      chat.InitialSystemMessage.Role,      // Defina o papel (ex: "user" ou "assistant")
			Content:   chat.InitialSystemMessage.Content,   // Conteúdo da mensagem inicial
			Tokens:    int32(chat.TokenUsage),              // Uso de tokens da mensagem
			CreatedAt: chat.InitialSystemMessage.CreatedAt, // Defina a data de criação
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepositoryMySQL) FindChatById(ctx context.Context, chatID string) (*entity.Chat, error) {
	chat := &entity.Chat{}
	res, err := r.Queries.FindChatByID(ctx, chatID)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	chat.ID = res.ID
	chat.UserID = res.UserID
	chat.Status = res.Status
	chat.TokenUsage = int(res.TokenUsage)
	chat.Config = &entity.ChatConfig{
		Model: &entity.Model{
			Name:      res.Model,
			MaxTokens: int(res.ModelMaxTokens),
		},
		Temperature:      float32(res.Temperature),
		TopP:             float32(res.TopP),
		N:                int(res.N),
		Stop:             []string{res.Stop},
		MaxTokens:        int(res.MaxTokens),
		PresencePenalty:  float32(res.PresencePenalty),
		FrequencyPenalty: float32(res.FrequencyPenalty),
	}

	messages, err := r.Queries.FindMessagesByChatID(ctx, chatID)

	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		chat.Messages = append(chat.Messages, &entity.Message{
			ID:        message.ID,
			Content:   message.Content,
			Role:      message.Role,
			Tokens:    int(message.Tokens),
			Model:     &entity.Model{Name: message.Model},
			CreatedAt: message.CreatedAt,
		})
	}

	erasedMessages, err := r.Queries.FindErasedMessagesByChatID(ctx, chatID)

	if err != nil {
		return nil, err
	}

	for _, message := range erasedMessages {
		chat.ErasedMessages = append(chat.ErasedMessages, &entity.Message{
			ID:        message.ID,
			Content:   message.Content,
			Role:      message.Role,
			Tokens:    int(message.Tokens),
			Model:     &entity.Model{Name: message.Model},
			CreatedAt: message.CreatedAt,
		})
	}

	return chat, nil

}

func (r *ChatRepositoryMySQL) SaveChat(ctx context.Context, chat *entity.Chat) error {
	params := db.SaveChatParams{
		ID:               chat.ID,
		UserID:           chat.UserID,
		InitialMessageID: chat.InitialSystemMessage.ID,          // Adicione o ID da mensagem inicial
		Status:           chat.Status,                           // Adicione o status do chat
		TokenUsage:       int32(chat.TokenUsage),                // Uso de tokens do chat
		Model:            chat.Config.Model.Name,                // Nome do modelo
		ModelMaxTokens:   int32(chat.Config.Model.MaxTokens),    // Máximo de tokens do modelo
		Temperature:      float64(chat.Config.Temperature),      // Temperatura do modelo
		TopP:             float64(chat.Config.TopP),             // Top P do modelo
		N:                int32(chat.Config.N),                  // N do modelo
		Stop:             chat.Config.Stop[0],                   // Stop do modelo
		MaxTokens:        int32(chat.Config.MaxTokens),          // Máximo de tokens
		PresencePenalty:  float64(chat.Config.PresencePenalty),  // Penalidade de presença
		FrequencyPenalty: float64(chat.Config.FrequencyPenalty), // Penalidade de frequência
		UpdatedAt:        time.Now(),
	}

	err := r.Queries.SaveChat(
		ctx,
		params,
	)

	if err != nil {
		return err
	}

	//deletar todas mensagens
	err = r.Queries.DeleteChatMessages(ctx, chat.ID)
	if err != nil {
		return err
	}

	//deletar mensagens apagadas
	err = r.Queries.DeleteErasedChatMessages(ctx, chat.ID)
	if err != nil {
		return err
	}

	i := 0

	// adicionar novas mensagens
	for _, message := range chat.Messages {
		err = r.Queries.AddMessage(
			ctx,
			db.AddMessageParams{
				ID:        message.ID,            // Adicione um ID único para a mensagem
				ChatID:    chat.ID,               // ID do chat associado
				Role:      message.Role,          // Defina o papel (ex: "user" ou "assistant")
				Content:   message.Content,       // Conteúdo da mensagem inicial
				Tokens:    int32(message.Tokens), // Uso de tokens da mensagem
				Model:     chat.Config.Model.Name,
				CreatedAt: message.CreatedAt, // Defina a data de criação
				OrderMsg:  int32(i),
				Erased:    false,
			},
		)
		if err != nil {
			return err
		}
		i++
	}

	//adicionar mensagens apagadas
	i = 0
	for _, message := range chat.ErasedMessages {
		err = r.Queries.AddMessage(
			ctx,
			db.AddMessageParams{
				ID:        message.ID,            // Adicione um ID único para a mensagem
				ChatID:    chat.ID,               // ID do chat associado
				Role:      message.Role,          // Defina o papel (ex: "user" ou "assistant")
				Content:   message.Content,       // Conteúdo da mensagem inicial
				Tokens:    int32(message.Tokens), // Uso de tokens da mensagem
				Model:     chat.Config.Model.Name,
				CreatedAt: message.CreatedAt, // Defina a data de criação
				OrderMsg:  int32(i),
				Erased:    true,
			},
		)
		if err != nil {
			return err
		}

		i++
	}

	return nil

}
