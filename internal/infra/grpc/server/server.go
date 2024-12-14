package server

import (
	"net"

	"github.com/gpt_chat/chat_service/internal/infra/grpc/pb"
	"github.com/gpt_chat/chat_service/internal/infra/grpc/service"
	"github.com/gpt_chat/chat_service/internal/usecase/chatcompletionstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	ChatComplementationStreamUseCase chatcompletionstream.ChatCompletionUseCase
	ChatConfigStream                 chatcompletionstream.ChatCompletionConfigInputDTO
	ChatService                      service.ChatService
	Port                             string
	AuthToken                        string
	StreamChannel                    chan chatcompletionstream.ChatCompletionOutputDTO
}

func NewGRPCServer(chatCompletionStreamUseCase chatcompletionstream.ChatCompletionUseCase, chatConfigStream chatcompletionstream.ChatCompletionConfigInputDTO, port string, authToken string, streamChannel chan chatcompletionstream.ChatCompletionOutputDTO) *GRPCServer {
	chatService := service.NewChatService(chatCompletionStreamUseCase, chatConfigStream, streamChannel)
	return &GRPCServer{
		ChatComplementationStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:                 chatConfigStream,
		Port:                             port,
		AuthToken:                        authToken,
		StreamChannel:                    streamChannel,
		ChatService:                      *chatService,
	}
}

func (g *GRPCServer) AuthInterceptor(service interface{}, serviceStream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	context := serviceStream.Context()
	md, ok := metadata.FromIncomingContext(context)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	if token[0] != g.AuthToken {
		return status.Error(codes.Unauthenticated, "authorization token is invalid")
	}

	return handler(service, serviceStream)
}

func (g *GRPCServer) Start() {
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(g.AuthInterceptor),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChatServiceServer(grpcServer, &g.ChatService)

	lis, err := net.Listen("tcp", ":"+g.Port)
	if err != nil {
		panic(err.Error())
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err.Error())
	}

}
