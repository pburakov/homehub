# Hub Agent

Reports Hub status to the main server and responds to user queries. 

## Prerequisites

- proto v3.6.0
- gRPC 1.13.1

Use Homebrew to install `protobuf` and `grpc` packages. Additionally:

```bash
$ go get -u google.golang.org/grpc
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

## Developing

Before developing, make sure to generate the latest schema definitions by running (make sure `agent` is the PWD):

```bash
$ mkdir schema 
$ protoc --go_out=plugins=grpc:schema ../schema/src/main/proto/*.proto --proto_path=../schema/src/main/proto
```