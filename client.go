package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/johnson7543/grpcChatServer/chatserver"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {

	// fmt.Println("Enter Server IP:Port ::: ")
	// reader := bufio.NewReader(os.Stdin)
	// serverID, err := reader.ReadString('\n')

	// if err != nil {
	// 	log.Printf("Failed to read from console :: %v", err)
	// }
	// serverID = strings.Trim(serverID, "\r\n")
	serverID := "localhost:5000"
	log.Info().Msgf("Connecting to : " + serverID)

	//connect to grpc server
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatal().Msgf("Faile to conncet to gRPC server :: %v", err)
	}
	defer conn.Close()

	//call ChatService to create a stream
	client := chatserver.NewServiceClient(conn)

	// create metadata with the client name
	clientName := clientConfig()
	md := metadata.New(map[string]string{"client-name": clientName})

	// create a context with the metadata
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := client.ChatService(ctx)
	if err != nil {
		log.Fatal().Msgf("Failed to call ChatService :: %v", err)
	}

	// implement communication with gRPC server
	ch := clienthandle{
		stream:     stream,
		clientName: clientName,
	}
	go ch.sendMessage()
	go ch.receiveMessage()

	//blocker
	bl := make(chan bool)
	<-bl

}

//clienthandle
type clienthandle struct {
	stream     chatserver.Service_ChatServiceClient
	clientName string
}

func clientConfig() string {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Error().Msgf(" Failed to read from console :: %v", err)
	}
	return strings.Trim(name, "\r\n")

}

func (ch *clienthandle) sendMessage() {

	for {

		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Error().Msgf(" Failed to read from console :: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Error().Msgf("Error while sending message to server :: %v", err)
		}

	}

}

func (ch *clienthandle) receiveMessage() {

	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Error().Msgf("Error in receiving message from server :: %v", err)
		}

		//print message to client's console
		fmt.Printf("%s : %s \n", mssg.Name, mssg.Body)

	}
}
