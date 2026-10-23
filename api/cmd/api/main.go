package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/K-T-Ng/ktblog/api/internal/httpapi"
	"github.com/K-T-Ng/ktblog/api/internal/posts"
)

func main() {
	cfg, err := loadConfig(os.Getenv)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := &http.Client{Timeout: 10 * time.Second}
	allowlist := posts.NewAllowlist(cfg.PostsURL, client)
	if err := allowlist.Refresh(ctx); err != nil {
		log.Fatalf("initial posts.json fetch: %v", err)
	}
	go allowlist.Run(ctx, cfg.RefreshInterval)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpapi.NewHandler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %q", srv.Addr)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("received shutdown signal")
	case err := <-errCh:
		log.Fatalf("unexpected error for launching server: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Printf("server stopped")
}
