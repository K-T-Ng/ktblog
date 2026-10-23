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
ADDR=127.0.0.1:8080                       # Address to listen on (host:port)
POSTS_URL=https://ktblog.dev/posts.json   # Where to fetch post slugs; use http://localhost:4321/posts.json for local dev
REFRESH_INTERVAL=10m                      # How often to refetch post slugs
```

## Test

```shell
go test ./...
```
