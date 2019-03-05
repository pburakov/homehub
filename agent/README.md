# Hub Agent

Reports Hub status to the main server and responds to user queries. 

## Generating proto sources (macOS)

Before developing, make sure to generate the latest schema definitions.

### Prerequisites

- Go 1.11.4+
- proto v3.6.0
- gRPC 1.13.1

Use Homebrew to install `go`, `protobuf` and `grpc` packages. Additionally, make sure `agent` is the CWD and run:

```bash
$ go get -u google.golang.org/grpc
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

When Go finished downloading these dependencies, run:

```bash
$ mkdir schema 
$ protoc --go_out=plugins=grpc:schema ../schema/src/main/proto/*.proto --proto_path=../schema/src/main/proto
```

Schema package can now be used in the code:

```go
import "io.pburakov/homehub/agent/schema"
``` 

## Build Application (macOS)

Make sure `agent` is the CWD. 

To build for Raspberry Pi 3 B+, run: 

```bash
$ env GOOS=linux GOARCH=arm GOARM=7 go build
```

Go will download required packages and install them into Go path. This will generate `agent` executable. 

Run with:

```bash
$ ./agent
```

Distribution package consists of generated `agent` executable and `conf` directory.

## Build & Run Application (Raspbian Linux)

### Prerequisites

- Go 1.11.4+
- motion 4.0

Use `apt-get` to install these packages. If `apt-get` does not contain the required Go version, download and unpack specific Go version:

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

Instructions for generating schema protobuf code for Raspbian Linux are to be added. Until then, schema sources are included with the repo. 

Make sure `agent` is the PWD. To build under Raspbian environment, simply run:  

```bash
$ go build
```

Run binary with:

```bash
$ ./agent
``` 
