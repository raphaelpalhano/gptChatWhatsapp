package service

import (
	"log"

	"github.com/gpt_chat/chat_service/internal/infra/grpc/pb"
	"github.com/gpt_chat/chat_service/internal/usecase/chatcompletionstream"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	ChatComplementationStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream                 chatcompletionstream.ChatCompletionConfigInputDTO
	StreamChannel                    chan chatcompletionstream.ChatCompletionOutputDTO
}

func NewChatService(chatComplementationStreamUseCase chatcompletionstream.ChatCompletionUseCase, chatConfigStream chatcompletionstream.ChatCompletionConfigInputDTO, streamChannel chan chatcompletionstream.ChatCompletionOutputDTO) *ChatService {
	return &ChatService{
		ChatComplementationStreamUseCase: chatComplementationStreamUseCase,
		ChatConfigStream:                 chatConfigStream,
		StreamChannel:                    streamChannel,
	}
}

func (c *ChatService) ChatStream(req *pb.ChatRequest, stream pb.ChatService_ChatStreamServer) error {
	chatConfig := chatcompletionstream.ChatCompletionConfigInputDTO{
		Temperature:          c.ChatConfigStream.Temperature,
		N:                    c.ChatConfigStream.N,
		TopP:                 c.ChatConfigStream.TopP,
		Stop:                 c.ChatConfigStream.Stop,
		MaxTokens:            c.ChatConfigStream.MaxTokens,
		Model:                c.ChatConfigStream.Model,
		InitialSystemMessage: c.ChatConfigStream.InitialSystemMessage,
		ModelMaxTokens:       c.ChatConfigStream.ModelMaxTokens,
	}

	input := chatcompletionstream.ChatCompletionInputDTO{
		UserMessage: req.GetUserMessage(),
		UserID:      req.GetUserId(),
		ChatID:      req.GetChatId(),
		Config:      chatConfig,
	}

	ctx := stream.Context()

	go func() {
		for msg := range c.StreamChannel {
			if err := stream.Send(&pb.ChatResponse{
				ChatId:  msg.ChatID,
				UserId:  msg.UserID,
				Content: msg.Content,
			}); err != nil {
				log.Printf("Erro ao enviar mensagem para o cliente: %v", err)
				return
			}
		}
	}()

	_, err := c.ChatComplementationStreamUseCase.Execute(ctx, input)
	if err != nil {
		return err
	}

	return nil

}
