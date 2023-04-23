# ChatServer impl using golang and gRPC (bidi)

## How to run the app

```bash
go run server.go

# split a new bash

go run client.go

localhost:5000
```

## Demo Result

### Server side

```bash
$ go run server.go
2023/04/23 21:37:54 Listening @ : 5000
2023/04/23 21:38:27 {johnson hi 19727887 498081}
2023/04/23 21:38:36 {angel hi 39984059 131847}
2023/04/23 21:38:38 {angel hi 11902081 131847}
2023/04/23 21:38:40 {johnson hi 74941318 498081}
```

### Client side

```bash

$ go run client.go
Enter Server IP:Port ::: 
localhost:5000
2023/04/23 21:38:22 Connecting : localhost:5000
Your Name : johnson
hi
angel : hi 
angel : hi 
hi
```

```bash

$ go run client.go
Enter Server IP:Port ::: 
localhost:5000
2023/04/23 21:38:33 Connecting : localhost:5000
Your Name : angel
johnson : hi 
hi
hi
johnson : hi 
```
