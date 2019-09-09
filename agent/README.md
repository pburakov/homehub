# Hub Agent

Reports Hub status to the main server and responds to user queries. 

## Prerequisites

- Go 1.13
- proto v3.6.0
- gRPC 1.13.1

Use Homebrew to install `go`, `protobuf` and `grpc` packages.

## Developing

Following instructions assume `agent` is the current working directory.

Make sure to set your `GOPATH` to `.gopath` before starting developing. 

This repo uses Maven for automatic sources generation. To generate latest schema sources, run Maven compilation:

```bash
mvn compile
```

Maven will create `GOPATH` directory `agent/.gopath`, download and install necessary dependency modules, then attempt a compilation. Please note that due to flakiness of Go dependency resolution, Maven is set to ignore errors of `go get` commands.

Once generated, schema files can be used in code: 

```go
import "io.pburakov/homehub/agent/schema"
```

## Build Application

Distributive binary packaging is managed by Maven. Run: 

```bash
mvn package
```

Maven wil build Raspberry Pi 3 B+ build, identical to running a command `env GOOS=linux GOARCH=arm GOARM=7 go build`.

Distribution package consists of generated `agent` executable and `conf` directory, found in `target/dist` directory.
