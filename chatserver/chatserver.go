package chatserver

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type message struct {
	ClientName       string
	Body             string
	UniqueCode       int
	ClientUniqueCode int
}

type MessageQueue struct {
	Messages []message
	mu       sync.Mutex
}

type clientConnection struct {
	clientStream Service_ChatServiceServer
	mu           sync.Mutex
}

var (
	messageQueue      = MessageQueue{}
	clientConnections = make(map[int]*clientConnection)
)

type ChatServer struct {
}

//define ChatService
func (is *ChatServer) ChatService(clientStream Service_ChatServiceServer) error {

	clientUniqueCode := rand.Intn(1e6)

	clientConn := &clientConnection{
		clientStream: clientStream,
	}

	clientConnections[clientUniqueCode] = clientConn

	errorChannel := make(chan error)

	// receive messages - init a go routine
	go receiveFromStream(clientConn, clientUniqueCode, errorChannel)

	// send messages - init a go routine
	go sendToStream(clientConn, clientUniqueCode, errorChannel)

	return <-errorChannel

}

//receive messages
func receiveFromStream(clientConn *clientConnection, clientUniqueCode_ int, errorChannel chan error) {

	//implement a loop
	for {
		rm, err := clientConn.clientStream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from client :: %v", err)
			errorChannel <- err
		} else {

			messageQueue.mu.Lock()

			messageQueue.Messages = append(messageQueue.Messages, message{
				ClientName:       rm.Name,
				Body:             rm.Body,
				UniqueCode:       rand.Intn(1e8),
				ClientUniqueCode: clientUniqueCode_,
			})

			log.Printf("%v", messageQueue.Messages[len(messageQueue.Messages)-1])

			messageQueue.mu.Unlock()
		}
	}
}

//send message
func sendToStream(clientConn *clientConnection, clientUniqueCode_ int, errorChannel chan error) {

	//implement a loop
	for {

		//loop through messages in Messages
		for {
			time.Sleep(500 * time.Millisecond)

			messageQueue.mu.Lock()

			if len(messageQueue.Messages) == 0 {
				messageQueue.mu.Unlock()
				break
			}

			senderUniqueCode := messageQueue.Messages[0].ClientUniqueCode
			senderName := messageQueue.Messages[0].ClientName
			clientMessage := messageQueue.Messages[0].Body

			messageQueue.mu.Unlock()

			// send message to all connected clients except the sender
			for clientUC, conn := range clientConnections {
				if clientUC != senderUniqueCode {

					conn.mu.Lock()

					err := conn.clientStream.Send(&FromServer{
						Name: senderName,
						Body: clientMessage,
					})

					conn.mu.Unlock()

					if err != nil {
						errorChannel <- err
					}
				}
			}

			messageQueue.mu.Lock()

			if len(messageQueue.Messages) > 1 {
				messageQueue.Messages = messageQueue.Messages[1:] // delete the message at index 0 after sending to receiver
			} else {
				messageQueue.Messages = []message{}
			}

			messageQueue.mu.Unlock()
		}

		time.Sleep(100 * time.Millisecond)
	}
}
