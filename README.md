# ChatServer impl using golang and gRPC (bidi)

## How to run the app

```bash
go run server.go

# split some new bash

go run client.go

# enter user name
# enter room name

```

## Demo Result

### Server side

```bash
$ go run server.go
2023/04/29 17:59:02 Listening @ : 5000
2023/04/29 17:59:42 Johnson is joining room: Cooper's Room
2023/04/29 17:59:42 There are 1 people in Cooper's Room now.
2023/04/29 17:59:51 Angel is joining room: Cooper's Room
2023/04/29 17:59:51 There are 2 people in Cooper's Room now.
2023/04/29 17:59:56 {Angel Cooper's Room Hi 27131847 727887}
2023/04/29 17:59:56 Sending message : Hi from Cooper's Room
2023/04/29 18:00:00 {Johnson Cooper's Room Hi 39984059 498081}
2023/04/29 18:00:00 Sending message : Hi from Cooper's Room
```

### Client side

```bash

$ go run client.go
2023/04/29 17:59:19 Connecting to : localhost:5000
Your Name : Johnson
SYSTEM : Welcome to the chat server! Please enter a room name: 
Cooper's Room          
Angel : Hi 
Hi
```

```bash

$ go run client.go
2023/04/29 17:59:20 Connecting to : localhost:5000
Your Name : Angel
SYSTEM : Welcome to the chat server! Please enter a room name: 
Cooper's Room
Hi
Johnson : Hi 
```
