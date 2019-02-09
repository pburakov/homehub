# Hub Agent

Reports Hub status to the main server and responds to user queries. 

## Generating proto sources (macOS)

Before developing, make sure to generate the latest schema definitions.

### Prerequisites

- Go 1.11.4+
- proto v3.6.0
- gRPC 1.13.1

Use Homebrew to install `go`, `protobuf` and `grpc` packages. Additionally:

```bash
$ go get -u google.golang.org/grpc
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

### Build Proto files

Make sure `agent` is the PWD:

```bash
$ mkdir schema 
$ protoc --go_out=plugins=grpc:schema ../schema/src/main/proto/*.proto --proto_path=../schema/src/main/proto
```

## Build Application

Make sure `agent` is the PWD. 

To build for Raspberry Pi 3 B+, run: 

```bash
$ env GOOS=linux GOARCH=arm GOARM=7 go build
```

Go will download required packages and install them into Go path. This will generate `homehub` executable. 

Run with:

```bash
$ ./homehub
``` 