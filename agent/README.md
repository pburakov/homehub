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

### Build

Make sure `agent` is the PWD:

```bash
$ mkdir schema 
$ protoc --go_out=plugins=grpc:schema ../schema/src/main/proto/*.proto --proto_path=../schema/src/main/proto
```

## Generating proto sources (Raspbian Linux)

Instructions for Raspbian Linux are to be added. Until then, proto sources are included with the repo. 

## Building and running under Raspbian Linux

### Prerequisites

- Go 1.11.4+
- ffmpeg

Download and unpack latest Go version:

```bash
$ wget https://storage.googleapis.com/golang/go1.11.4.linux-armv6l.tar.gz
$ sudo tar -C /usr/local -xvf go1.11.4.linux-armv6l.tar.gz
```

Verify installed version (should match 1.11.4):

```bash
$ go version
```

Define `GOPATH` and path to Go binaries:

```bash
$ cat >> ~/.bashrc << 'EOF'
  export GOPATH=$HOME/go
  export PATH=/usr/local/go/bin:$PATH:$GOPATH/bin
  EOF
```

Use `apt-get` to install `ffmpeg` package (`sudo apt-get install ffmpeg`).

### Build

Make sure `agent` is the PWD:

```bash
$ go build
```

Go will download required packages and install them into Go path.  

### Run

This will generate `homehub` executable. Run with:

```bash
$ ./homehub
``` 