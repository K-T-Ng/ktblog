# ktblog API

The blog's backend: a Go HTTP API for per-post views and likes.

## Prerequisites

Go 1.27 or later.

## Run locally

```shell
go run ./cmd/api
```

Optional environment variables:
```shell
ADDR=127.0.0.1:8080  # Address to listen on (host:port)
```

## Test

```shell
go test ./...
```
