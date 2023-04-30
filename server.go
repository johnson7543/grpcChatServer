package main

import (
	"net"
	"os"

	"github.com/johnson7543/grpcChatServer/chatserver"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
)

func main() {

	//assign port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000" //default Port set to 5000 if PORT is not set in env
	}

	//init listener
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatal().Msgf("Could not listen @ %v :: %v", Port, err)
	}
	log.Info().Msgf("Listening @ : " + Port)

	//gRPC server instance
	grpcserver := grpc.NewServer()

	//register ChatService
	cs := chatserver.ChatServer{}
	chatserver.RegisterServiceServer(grpcserver, &cs)

	//grpc listen and serve
	err = grpcserver.Serve(listen)
	if err != nil {
		log.Fatal().Msgf("Failed to start gRPC Server :: %v", err)
	}

}
