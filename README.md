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
{"level":"info","time":"2023-04-30T18:43:24+08:00","message":"Listening @ : 5000"}
{"level":"info","time":"2023-04-30T18:43:55+08:00","message":"Johnson is joining room: Cooper's Room"}
{"level":"info","time":"2023-04-30T18:43:55+08:00","message":"There are 1 people in Cooper's Room now."}
{"level":"info","time":"2023-04-30T18:44:01+08:00","message":"Angel is joining room: Cooper's Room"}
{"level":"info","time":"2023-04-30T18:44:01+08:00","message":"There are 2 people in Cooper's Room now."}
{"level":"info","ClientName":"Johnson","Room":"Cooper's Room","Body":"hi","UniqueCode":27131847,"ClientUniqueCode":498081,"time":"2023-04-30T18:44:05+08:00","message":"Receiving message : hi from Johnson"}
{"level":"info","Name":"Johnson","Body":"hi","time":"2023-04-30T18:44:05+08:00","message":"Sending message : hi from Cooper's Room"}
{"level":"info","ClientName":"Angel","Room":"Cooper's Room","Body":"hi","UniqueCode":39984059,"ClientUniqueCode":727887,"time":"2023-04-30T18:44:08+08:00","message":"Receiving message : hi from Angel"}
{"level":"info","Name":"Angel","Body":"hi","time":"2023-04-30T18:44:08+08:00","message":"Sending message : hi from Cooper's Room"}
```

### Client side

```bash
$ go run client.go
{"level":"info","time":"2023-04-30T18:43:26+08:00","message":"Connecting to : localhost:5000"}
Your Name : Johnson     
SYSTEM : Welcome to the chat server! Please enter a room name: 
Cooper's Room
hi
Angel : hi 
```

```bash
$ go run client.go
{"level":"info","time":"2023-04-30T18:43:28+08:00","message":"Connecting to : localhost:5000"}
Your Name : Angel
SYSTEM : Welcome to the chat server! Please enter a room name: 
Cooper's Room
Johnson : hi 
hi
```
