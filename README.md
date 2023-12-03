### Directory Structure

- `api`: A basic REST gateway, forwarding requests onto service(s).
- `racing`: A very bare-bones racing service.
- `sports`: A very bare-bones sports service.

```
grpc-gateway/
├─ api/
│  ├─ proto/
│  ├─ main.go
├─ racing/
│  ├─ db/
│  ├─ proto/
│  ├─ service/
│  ├─ main.go
├─ sports/
│  ├─ db/
│  ├─ proto/
│  ├─ service/
│  ├─ main.go
├─ README.md
```

### Getting Started

1. Install Go (latest).

```bash
brew install go
```

... or [see here](https://golang.org/doc/install).

2. Install `protoc`

```
brew install protobuf
```

... or [see here](https://grpc.io/docs/protoc-installation/).

2. In a terminal window, start our racing service...

```bash
cd ./racing

go build && ./racing
➜ INFO[0000] gRPC server listening on: localhost:9000
```

3. In a terminal window, start our racing service...

```bash
cd ./sports

go build && ./sports
➜ INFO[0000] gRPC server listening on: localhost:9001
```

4. In another terminal window, start our api service...

```bash
cd ./api

go build && ./api
➜ INFO[0000] API server listening on: localhost:8000
```

5. Make a request for races... 

```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d $'{
  "filter": {}
}'

curl --location --request GET 'http://localhost:8000/v1/race?id=1' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "id": 2
}'

curl --location --request POST 'http://localhost:8000/v1/list-sports' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "filter": {
        "visible":true,
        "column":"result",
        "order_by": "desc"
    }
}'

curl --location --request POST 'http://localhost:8000/v1/list-races' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "filter": {
        "visible":true,
        "column":"advertised_start_time",
        "order_by": "asc"
    }
}'
```

**Note:**

To aid in proto generation following any changes, you can run `go generate ./...` from `api` and `racing` and `sports` directories.

Before you do so, please ensure you have the following installed. You can simply run the following command below in each of `api` and `racing` and `sports` directories.

```
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 google.golang.org/genproto/googleapis/api google.golang.org/grpc/cmd/protoc-gen-go-grpc google.golang.org/protobuf/cmd/protoc-gen-go
```
