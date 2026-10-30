# ktblog API

The blog's backend: a Go HTTP API for per-post views and likes.

## Prerequisites

Go 1.27 or later.

## Run locally

```shell
go run ./cmd/api
```

All configuration is through environment variables, and all are optional:

| Variable           | Default                         | Purpose                                                                                                                   |
|--------------------|---------------------------------|---------------------------------------------------------------------------------------------------------------------------|
| `ADDR`             | `127.0.0.1:8080`                | Address to listen on (host:port)                                                                                          |
| `ALLOWED_ORIGINS`  | `http://localhost:4321`         | Comma-separated origins allowed to call the API from a browser (CORS). Each is `scheme://host[:port]`, no trailing slash. |
| `POSTS_URL`        | `https://ktblog.dev/posts.json` | Where to fetch post slugs                                                                                                 |
| `REFRESH_INTERVAL` | `10m`                           | How often to refetch post slugs                                                                                           |

For local development with the site running (`npm run dev`), set `POSTS_URL=http://localhost:4321/posts.json`.

## Test

```shell
go test ./...
```
