package main

import (
	"database/sql"
	"fmt"

	"github.com/gpt_chat/chat_service/configs"
	"github.com/gpt_chat/chat_service/internal/infra/grpc/server"
	"github.com/gpt_chat/chat_service/internal/infra/repository"
	"github.com/gpt_chat/chat_service/internal/infra/web"
	"github.com/gpt_chat/chat_service/internal/infra/web/webserver"
	"github.com/gpt_chat/chat_service/internal/usecase/chatcompletionstream"
	chatcompletion "github.com/gpt_chat/chat_service/internal/usecase/completion"
	"github.com/sashabaranov/go-openai"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	conn, err := sql.Open(configs.DBDriber, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	repo := repository.NewChatRepositoryMySQL(conn)
	client := openai.NewClient(configs.OpenAIApiKey)

	chatConfig := chatcompletion.ChatCompletionConfigInputDTO{
		Model:                configs.Model,
		ModelMaxTokens:       configs.ModelMaxTokens,
		Temperature:          float32(configs.Temperature),
		TopP:                 float32(configs.TopP),
		N:                    configs.N,
		Stop:                 configs.Stop,
		MaxTokens:            configs.MaxTokens,
		InitialSystemMessage: configs.InitialChatMessage,
	}

	chatConfigStream := chatcompletionstream.ChatCompletionConfigInputDTO{
		Model:                configs.Model,
		ModelMaxTokens:       configs.ModelMaxTokens,
		Temperature:          float32(configs.Temperature),
		TopP:                 float32(configs.TopP),
		N:                    configs.N,
		Stop:                 configs.Stop,
		MaxTokens:            configs.MaxTokens,
		InitialSystemMessage: configs.InitialChatMessage,
	}

	usecase := chatcompletion.NewChatCompletionUseCase(repo, client)
	streamChanel := make(chan chatcompletionstream.ChatCompletionOutputDTO)
	usecaseStream := chatcompletionstream.NewChatCompletionUseCase(repo, client, streamChanel)

	fmt.Println("Server GRPC starting on port ", configs.GRPCServerPort)
	grpcServer := server.NewGRPCServer(*usecaseStream, chatConfigStream, configs.GRPCServerPort, configs.AuthToken, streamChanel)
	go grpcServer.Start()

	webserver := webserver.NewWebServer(":" + configs.WebServerPort)
	webserverChatHandler := web.NewWebChatGPTHandler(*usecase, chatConfig, configs.AuthToken)
	webserver.AddHandler("/chat", webserverChatHandler.Handle)

	fmt.Println("Server Runing on port ", configs.WebServerPort)
	webserver.Start()
}
