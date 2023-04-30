package chatserver

import (
	fmt "fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/metadata"
)

type message struct {
	ClientName       string
	Room             string
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

type Room struct {
	name         string
	clients      map[int]*clientConnection // clientUniqueCode -> *clientConnection
	clientsMutex sync.RWMutex
}

var (
	messageQueue      = MessageQueue{}
	clientConnections = make(map[int]*clientConnection) // all connected clients
	rooms             = make(map[string]*Room)          // rooms to client name mapping
	roomsMutex        = sync.RWMutex{}
)

type ChatServer struct {
}

//define ChatService
func (cs *ChatServer) ChatService(clientStream Service_ChatServiceServer) error {
	// Access the context from the stream
	ctx := clientStream.Context()
	errorChannel := make(chan error)

	// Get the client's name from the context metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err := fmt.Errorf("failed to get metadata from context")
		errorChannel <- err
		return <-errorChannel
	}
	clientName := md["client-name"][0]

	clientUniqueCode := rand.Intn(1e6)
	clientConn := &clientConnection{
		clientStream: clientStream,
	}
	clientConnections[clientUniqueCode] = clientConn

	// prompt the client to join a room
	if err := clientConn.clientStream.Send(&FromServer{
		Name: "SYSTEM",
		Body: "Welcome to the chat server! Please enter a room name:",
	}); err != nil {
		log.Printf("Error sending message to client %d: %v", clientUniqueCode, err)
		delete(clientConnections, clientUniqueCode)
		errorChannel <- err
		return <-errorChannel
	}

	// receive room name from the client
	joinRequest, err := clientConn.clientStream.Recv()
	if err != nil {
		log.Error().
			Err(err).Int("clientUniqueCode", clientUniqueCode).Msgf("Error receiving room name from client")
		// log.Printf("Error receiving room name from client %d: %v", clientUniqueCode, err)
		delete(clientConnections, clientUniqueCode)
		errorChannel <- err
		return <-errorChannel
	}

	log.Info().Msgf("%v is joining room: %v", clientName, joinRequest.Body)
	// log.Printf("%v is joining room: %v", clientName, joinRequest.Body)

	roomName := joinRequest.Body

	// add the client to the room
	roomsMutex.Lock()
	room, ok := rooms[roomName]
	if !ok {
		room = &Room{
			name:    roomName,
			clients: make(map[int]*clientConnection),
		}
		rooms[roomName] = room
	}
	room.AddClientToRoom(clientConn, clientUniqueCode)
	roomsMutex.Unlock()

	// receive messages - init a go routine
	go receiveFromStream(clientConn, clientUniqueCode, errorChannel)

	// send messages - init a go routine
	go sendToStream(clientConn, clientUniqueCode, errorChannel)

	return <-errorChannel
}

func getRoomForClient(clientUniqueCode_ int) (*Room, error) {
	for _, room := range rooms {
		room.clientsMutex.RLock()
		_, ok := room.clients[clientUniqueCode_]
		room.clientsMutex.RUnlock()

		if ok {
			return room, nil
		}
	}

	return nil, fmt.Errorf("no room found for client %d", clientUniqueCode_)
}

func receiveFromStream(clientConn *clientConnection, clientUniqueCode_ int, errorChannel chan error) {

	for {
		receiveMessage, err := clientConn.clientStream.Recv()
		if err != nil {
			log.Error().Err(err).Msgf("Error in receiving message from client")
			// log.Printf("Error in receiving message from client :: %v", err)
			errorChannel <- err
		} else {

			messageQueue.mu.Lock()
			roomsMutex.RLock()

			room, err := getRoomForClient(clientUniqueCode_)
			if err != nil {
				log.Error().Err(err).Msgf("Error receiving message from client %v in room %v", receiveMessage.Name, room.name)
				// log.Printf("Error receiving message from client %v in room %v: %v", receiveMessage.Name, room.name, err)
				errorChannel <- err
				roomsMutex.RUnlock()
				messageQueue.mu.Unlock()
			} else {

				uniqueCode := rand.Intn(1e8)

				log.Info().
					Str("ClientName", receiveMessage.Name).
					Str("Room", room.name).
					Str("Body", receiveMessage.Body).
					Int("UniqueCode", uniqueCode).
					Int("ClientUniqueCode", clientUniqueCode_).
					Msgf("Receiving message : %v from %v", receiveMessage.Body, receiveMessage.Name)

				messageQueue.Messages = append(messageQueue.Messages, message{
					ClientName:       receiveMessage.Name,
					Room:             room.name,
					Body:             receiveMessage.Body,
					UniqueCode:       uniqueCode,
					ClientUniqueCode: clientUniqueCode_,
				})

				roomsMutex.RUnlock()
				messageQueue.mu.Unlock()
			}
		}
	}
}

func sendToStream(clientConn *clientConnection, clientUniqueCode_ int, errorChannel chan error) {

	for {
		time.Sleep(500 * time.Millisecond)

		messageQueue.mu.Lock()

		if len(messageQueue.Messages) == 0 {
			messageQueue.mu.Unlock()
			continue
		}

		message := messageQueue.Messages[0]
		messageQueue.Messages = messageQueue.Messages[1:]

		messageQueue.mu.Unlock()

		room, ok := rooms[message.Room]
		if ok {
			err := room.Broadcast(message)

			if err != nil {
				log.Error().Msgf("Error broadcasting message to client %d in room %v: %v", message.ClientUniqueCode, room.name, err)
				// log.Printf("Error broadcasting message to client %d in room %v: %v", message.ClientUniqueCode, room.name, err)
				errorChannel <- err
			}
		} else {
			err := fmt.Errorf("can not find room for %v", message.ClientName)
			log.Error().Err(err).Msgf("Error broadcasting message to client %d in room %v", message.ClientUniqueCode, room.name)
			// log.Printf("Error broadcasting message to client %d in room %v: %v", message.ClientUniqueCode, room.name, err)
			errorChannel <- err
		}

	}
}

func (room *Room) AddClientToRoom(client *clientConnection, clientUniqueCode_ int) {
	room.clientsMutex.Lock()
	defer room.clientsMutex.Unlock()

	room.clients[clientUniqueCode_] = client

	log.Info().Msgf("There are %d people in %v now.", len(room.clients), room.name)
	// log.Printf("There are %d people in %v now.", len(room.clients), room.name)
}

func (room *Room) RemoveClientFromRoom(client *clientConnection, clientUniqueCode_ int) {
	room.clientsMutex.Lock()
	defer room.clientsMutex.Unlock()

	delete(room.clients, clientUniqueCode_)
}

func (room *Room) Broadcast(msg message) error {
	room.clientsMutex.RLock()
	defer room.clientsMutex.RUnlock()

	for clientUC, conn := range room.clients {
		if clientUC != msg.ClientUniqueCode {

			conn.mu.Lock()

			log.Info().
				Str("Name", msg.ClientName).
				Str("Body", msg.Body).
				Msgf("Sending message : %v from %v", msg.Body, room.name)
			// log.Printf("Sending message : %v from %v", msg.Body, room.name)

			err := conn.clientStream.Send(&FromServer{
				Name: msg.ClientName,
				Body: msg.Body,
			})

			conn.mu.Unlock()

			if err != nil {
				log.Error().Err(err).Msgf("Error broadcasting message to client %d in room %v", clientUC, room.name)
				// log.Printf("Error broadcasting message to client %d in room %v: %v", clientUC, room.name, err)
				return err
			}
		}
	}

	return nil
}
